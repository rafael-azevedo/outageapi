package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-oci8"
	"github.com/rafael-azevedo/outageapi/utils"
	"github.com/spf13/viper"
)

type OracleDB struct {
	Username    string
	Password    string
	HostName    string
	Port        string
	ServiceName string
}

type NodeStatus struct {
	NodeName string
	InOutage bool
}

type MultiStatus []NodeStatus

func (m *MultiStatus) GetNodeStatus(db *sql.DB) error {
	query := `SELECT NN.NODE_NAME, NVL((SELECT CASE WHEN NODE_GROUP_ID = 'c44dc002-2c7c-71e5-07b9-0a04a4530000' THEN 'true' ELSE 'false' END FROM OPC_NODES_IN_GROUP WHERE NODE_GROUP_ID = 'c44dc002-2c7c-71e5-07b9-0a04a4530000' AND NODE_ID = NN.NODE_ID), 'false') AS IN_OUTAGE FROM OPC_NODE_NAMES NN LEFT JOIN OPC_NODES N ON N.NODE_ID = NN.NODE_ID WHERE NN.NETWORK_TYPE = 1 AND NVL(N.NODE_ID,'BAD') <> 'BAD'`

	rows, err := db.Query(query)

	if err != nil {
		return err

	}

	defer rows.Close()

	for rows.Next() {
		var n NodeStatus
		err := rows.Scan(&n.NodeName, &n.InOutage)
		if err != nil {
			return err

		}
		*m = append(*m, n)
	}

	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func (m *MultiStatus) OutageStatus() error {

	env := utils.EnvKeys{}

	//Get configuratoin info from config file
	cfg := viper.New()
	cfg.SetConfigName("app")
	config := os.Getenv("OUTAGECONF")
	cfg.AddConfigPath(config)

	err := cfg.ReadInConfig()
	if err != nil {
		err = fmt.Errorf("Fatal error config file: %s \n", err)
		return err
	}

	oDB, err := utils.NewOracleDB(cfg, &env)
	if err != nil {
		err = fmt.Errorf("Could not create Oracle object from config file")
		return err
	}

	//Connect to the database
	ConnString := oDB.Username + "/" + oDB.Password + "@" + oDB.HostName + ":" + oDB.Port + "/" + oDB.ServiceName
	db, err := sql.Open("oci8", ConnString)
	if err != nil {
		return err
	}

	//Execute Query that builds slice of nodeStatus

	err = m.GetNodeStatus(db)
	if err != nil {
		return err
	}
	db.Close()
	return nil
}

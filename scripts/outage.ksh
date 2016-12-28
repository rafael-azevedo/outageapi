#! /usr/bin/ksh

ACTION="$1"

NODE="$2"

NOTES="$3"

SCRIPTS=/var/opt/OV/SPLS_scripts

if [ ${#} -ne 3 ]
   then
   echo "USAGE:"
   echo "  Required values in order of input:"
   echo "     Action to take:   ADD   or  REMOVE"
   echo "     FQDN - the fully qualified node name"
   echo "     Comments in quotes, Change and when can be removed"
   echo "       if decommission, other."
   echo ""
   exit
fi

case $NOTES in
"" )
  echo "Please provide a Change Ticket and comments"
  echo "  Such as temporarily in outage util DATE"
  echo ""
  echo "exiting"
  exit
;;
esac

OS=`uname -s`
case $OS in
"HP-UX" )

     OVPATH=/opt/OV/bin/OpC/utils
     OVBIN=/opt/OV/bin
     ;;
* )
#LINUX
     OVPATH=/opt/OV/bin/OpC/utils
     OVBIN=/opt/OV/bin
     ;;
esac

OPCNODE=$OVPATH/opcnode

LOGf=`find /var/opt/OV/log/node_add_logs -name "$NODE"`
echo "debug $LOGf"

#if [ ! -file $LOGf ]
LOG=$LOGf
if [ "$LOGf" == "" ]
  then

echo "debug2 $LOGf"
    LOG="/var/opt/OV/log/node_add_logs/$NODE"
fi

echo "debug3 LOG $LOG"

echo "************************************* | tee -a $LOG
echo "Processing OUTAGE Request" | tee -a $LOG
echo "*************************************

DATE_TIME=`date`
echo "$DATE_TIME" | tee -a $LOG

#
#Function call to validate if node in OVO DB and also echo usage if no node provided.
#

validate_node () {
    NodeName=$1
    if [ "$NodeName" = "" ]
       then
        echo "" | tee -a $LOG
	echo " Missing Parameter" | tee -a $LOG
	echo " " | tee -a $LOG
        echo " You need to specify a node name." | tee -a $LOG
        echo " usage:  outage_set.ksh <ACTION> <NODENAME>" | tee -a $LOG
        echo "" | tee -a $LOG
       else
        echo "Validating node $NodeName" | tee -a $LOG
        echo " " | tee -a $LOG

        $OPCNODE -list_id node_name=$NodeName 2>&1 > /dev/null
        result=$?

        if [ $result -gt 0 ]
           then
            echo "   ERROR" | tee -a $LOG
            echo "   Node not found in DB, not a valide node name" | tee -a $LOG
	    VALID="NO"
           else
            echo "   Valide node name" | tee -a $LOG
            echo "" | tee -a $LOG
	    VALID="YES"
           fi
       fi
}

#
# Begin program MAIN
#
#Process node add to outage group
#
case $ACTION in
  "ADD" )
     validate_node $NODE
     if [ "$VALID" = "YES" ]
        then
            echo "" | tee -a $LOG
            echo "   Adding node $NODE to outage group" | tee -a $LOG
            $OPCNODE -assign_node group_name=outage node_name=$NODE net_type="NETWORK_IP" | tee -a $LOG

            result=$?
            if [ $result -gt 0 ]
               then
                echo " " | tee -a $LOG
	        echo "   ERROR" | tee -a $LOG
	        echo "   Node $NODE not added to outage." | tee -a $LOG
	        echo "   Please contact the Tools Design Team for assistance." | tee -a $LOG
	        echo "   This is not a critical issue and does not require escalation." | tee -a $LOG
	        echo "   Please open a incident for assistance." | tee -a $LOG
                echo "ADDED : `date` : with errors : $NOTES." >> /var/opt/OV/log/OUTAGElog/$NODE
	       else
	        echo " " | tee -a $LOG
	        echo " Node $NODE successfully set to outage state." | tee -a $LOG
	        echo " " | tee -a $LOG
                echo "ADDED : `date` : $NOTES." >> /var/opt/OV/log/OUTAGElog/$NODE
               fi
            # This command will disable the node and prevent additional messages from entring the queue
            echo " " | tee -a $LOG
            echo "Disabling agent $NODE" | tee -a $LOG
             echo "$OPCNODE -chg_nodetype node_name=$NODE node_type=DISABLED net_type=NETWORK_IP"
            $OPCNODE -chg_nodetype node_name="$NODE" node_type=DISABLED net_type=NETWORK_IP | tee -a $LOG
            result=$?
		echo "result $result"
            if [ $result -ne 0 ]
               then
                cho " " | tee -a $LOG
	        echo "   ERROR" | tee -a $LOG
	        echo "   Node $NODE not set to DISABLED." | tee -a $LOG
	       else
	        echo " " | tee -a $LOG
	        echo " Node $NODE successfully set to DISABLED." | tee -a $LOG
	        echo " " | tee -a $LOG
               fi


        fi
    ;;
#
#Process node remove from outage group
#
  "REMOVE" )
     validate_node $NODE
     if [ "$VALID" = "YES" ]
        then
          # This command will ENABLE the node and prevent additional messages from entring the queue
            echo " " | tee -a $LOG
            echo "ENABLing agent $NODE" | tee -a $LOG
            $OPCNODE -chg_nodetype node_name="$NODE" node_type=CONTROLLED net_type=NETWORK_IP | tee -a $LOG
            result=$?
            if [ $result -gt 0 ]
              then
               echo " " | tee -a $LOG
               echo "   ERROR" | tee -a $LOG
               echo "   Node $NODE not set to CONTROLLED." | tee -a $LOG
              else
               echo " " | tee -a $LOG
               echo " Node $NODE successfully set to CONTROLLED." | tee -a $LOG
               echo " " | tee -a $LOG
            fi

            echo "removing node $NODE from outage group" | tee -a $LOG
            echo " " | tee -a $LOG
            $OPCNODE -deassign_node group_name=outage node_name=$NODE net_type=NETWORK_IP
            result=$?
            if [ $result -gt 0 ]
               then
                echo " " | tee -a $LOG
                echo "   ERROR" | tee -a $LOG
                echo "   Node $NODE not removed from outage state." | tee -a $LOG
                echo "   Please contact the Tools Design Team for assistance." | tee -a $LOG
                echo "   This is not a critical issue and does not require escalation." | tee -a $LOG
                echo "   Please open a incident for assistance." | tee -a $LOG
               else
                echo " " | tee -a $LOG
                echo " Node $NODE successfully removed from outage state." | tee -a $LOG
                echo " " | tee -a $LOG
             fi

             echo "opcsw -installed $NODE" | tee -a $LOG
             opcsw -installed $NODE

             $OVBIN/opcragt -cleanstart $NODE
             echo "brief pause required"
             sleep 30

             $OVBIN/opcragt -distrib -force $NODE
             echo "brief pause required"
	     sleep 30

             $OVBIN/opcragt -status $NODE | tee -a $LOG

             $SCRIPTS/SetUpHB.ksh  $NODE
             print "REMOVED : `date` : $NOTES." >> /var/opt/OV/log/OUTAGElog/$NODE

          #rm -f /var/opt/OV/log/OUTAGElog/$NODE
       fi

   ;;
  *)
    echo " "
    echo " Missing Parameters"
    echo " "
    echo " Choices are: "
    echo ""
    echo "      ADD  a node to outage"
    echo "      REMOVE  a node from outage"
    echo ""
    echo " ADDING a node to outage suspends monitoring and all OM alerts"
    echo " REMOVING a node from outage restores monitoring and all OM alters"
    echo ""
    echo " usage:  outage_set.ksh <ACTION> <NODENAME>"
    echo ""
   ;;
esac
echo "log $LOG"
echo "******* SET agent tracing off" >> $LOG
echo "SET agent tracing off" >> $LOG
echo "" >> $LOG
cmd = "/opt/OV/bin/ovdeploy -cmd '/opt/OV/support/ovtrccfg -off' -host $NodeName"

echo "$cmd" | tee -a $LOG

return = `$cmd 2>&1`

/opt/OV/bin/ovdeploy -cmd "/opt/OV/support/ovtrccfg -off" -host $NodeName

if [ $? != 0 ]
  then
        echo "    ***  Agent tracing off change failed. return = $return***" >> $LOG
  else
        echo " Agent tracing off configuration change successful: $return " >> $LOG
  fi

exit
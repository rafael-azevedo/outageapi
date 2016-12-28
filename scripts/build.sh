#!/bin/bash

#http://www.apache.org/licenses/LICENSE-2.0.txt
#
#
#
#Licensed under the Apache License, Version 2.0 (the "License");
#you may not use this file except in compliance with the License.
#You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
#Unless required by applicable law or agreed to in writing, software
#distributed under the License is distributed on an "AS IS" BASIS,
#WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#See the License for the specific language governing permissions and
#limitations under the License.

#Setting Bash to complain on errors
set -e
set -u
set -o pipefail

#Setting output format 
_fmt () {
  local color_debug="\x1b[35m"
  local color_info="\x1b[32m"
  local color_notice="\x1b[34m"
  local color_warning="\x1b[33m"
  local color_error="\x1b[31m"
  local colorvar=color_$1

  local color="${!colorvar:-$color_error}"
  local color_reset="\x1b[0m"
  echo -e "$(date -u +"%Y-%m-%d %H:%M:%S UTC") ${color}$(printf "[%s]" "${1}")${color_reset}";
}
#Defining variables necessary
NO_COLOR="${NO_COLOR:-}"
go_build=(go build -i)


#Defining working directories
__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__proj_dir="$(dirname "$__dir")"
_build_path="${__proj_dir}/bin/"

#Useage of script
_usage () {
  echo $"Usage: $0 {build|clean} {all|encryptcli|outageapi}"
  echo $"To see this massage again : $0 {-h|help}"
}


#Printing relevant information
_info (){ echo "$(_fmt info) ${*}" 1>&2 || true; }
_error (){ echo "$(_fmt error) ${*}" 1>&2 || true; }

#Building Outage API Binary 
_build_outageapi () {
_info "project path: ${__proj_dir}"
_info "Building outageapi at ${_build_path}"
mkdir -p "${_build_path}"
(cd "${__proj_dir}/cmd" && "${go_build[@]}" -o "${_build_path}/outageapi" . || exit 1)
}

#Build EncryptCli
_build_encryptcli (){
_info "project path: ${__proj_dir}"
_info "Building encryptCLI at ${_build_path}"
mkdir -p "${_build_path}"
(cd "${__proj_dir}/encryptCLI" && "${go_build[@]}" -o "${_build_path}/encryptCLI" . || exit 1)
}

#Build all binarys
_build_all (){
  _info "project path: ${__proj_dir}"
  _info "Building all at ${_build_path}"
  mkdir -p "${_build_path}"
  (cd "${__proj_dir}/cmd" && "${go_build[@]}" -o "${_build_path}/outageapi" . || exit 1)
  (cd "${__proj_dir}/encryptCLI" && "${go_build[@]}" -o "${_build_path}/encryptCLI" . || exit 1)
}

#Cleaning commands
_clean_outageapi (){
   _info "Cleaning ${_build_path}/outageapi"
   if [ -f ${_build_path}/outageapi ]; then
   rm ${_build_path}/outageapi
   fi
}

_clean_encryptcli (){
   _info "Cleaning ${_build_path}/encryptCLI"
   if [ -f ${_build_path}/encryptCLI ]; then
   rm ${_build_path}/encryptCLI
   fi
}

_clean_all (){
  _info "Cleaning ${_build_path}/outageapi and ${_build_path}/encryptCLI"
   if [ -f ${_build_path}/outageapi ]; then
   rm ${_build_path}/outageapi
   fi  
   
   if [ -f ${_build_path}/encryptCLI ]; then
   rm ${_build_path}/encryptCLI
   fi
}

#Main

#  Get build Variables
if  [ $# -lt 2 -o $# -gt 2 ]; then
    _usage
    exit 1
fi


case ${1} in 
  help|-h) 
    _usage
    exit 0
    ;;
  build|clean)
    action=$1
    shift
      case $1 in
        all|encryptcli|outageapi)
          binary=$1
          shift
          ;;
        *)
        _usage
        exit 1
        ;;
      esac
    ;;
  *)
    _usage
    exit 1
    ;;
esac

#build according to variables
_${action}_${binary}
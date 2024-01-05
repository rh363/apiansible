package main

/*
NAME=APIANSIBLE
AUTHOR=RH363
DATE=01/2024
COMPANY=SEEWEB
VERSION=1.1
*/
import (
	"errors"
	"net/http"

	ansible "github.com/febrianrendak/go-ansible"
	"github.com/gin-gonic/gin"
)

// RESPONSE MESSAGES
var RES_installed string = "[SUCCESS] INSTALLED"
var RES_uninstalled string = "[SUCCESS] UNINSTALLED"
var RES_added string = "[SUCCESS] ADDED"
var RES_removed string = "[SUCCESS] REMOVED"

// ERRORS:

// JSON ERROR

var ERR_bad_json_format error = errors.New("[ERROR] BAD JSON FORMAT")

// NFS PLAYBOOK ERROR

var ERR_cant_run_add_nfs_access_playbook error = errors.New("[ERROR] CANT RUN ADD NFS ACCESS PLAYBOOK")
var ERR_cant_run_remove_nfs_access_playbook error = errors.New("[ERROR] CANT RUN REMOVE NFS ACCESS PLAYBOOK")
var ERR_cant_run_install_nfs_playbook error = errors.New("[ERROR] CANT RUN INSTALL NFS PLAYBOOK")
var ERR_cant_run_uninstall_nfs_playbook error = errors.New("[ERROR] CANT RUN UNINSTALL NFS PLAYBOOK")

// SMB PLAYBOOK ERROR

var ERR_cant_run_add_smb_access_playbook error = errors.New("[ERROR] CANT RUN ADD SMB ACCESS PLAYBOOK")
var ERR_cant_run_remove_smb_access_playbook error = errors.New("[ERROR] CANT RUN REMOVE SMB ACCESS PLAYBOOK")
var ERR_cant_run_install_smb_playbook error = errors.New("[ERROR] CANT RUN INSTALL SMB PLAYBOOK")
var ERR_cant_run_uninstall_smb_playbook error = errors.New("[ERROR] CANT RUN UNINSTALL SMB PLAYBOOK")

// FILES PATH
// PLAYBOOK DIRECTORY PATH

var playbook_dir string = "/etc/ansible/playbooks/"

// NFS PLAYBOOK PATH

var add_nfs_access_playbook string = playbook_dir + "add-nfs-playbook.yml"
var remove_nfs_access_playbook string = playbook_dir + "remove-nfs-playbook.yml"
var install_nfs_playbook string = playbook_dir + "install-nfs-playbook.yml"
var uninstall_nfs_playbook string = playbook_dir + "uninstall-nfs-playbook.yml"

// SMB PLAYBOOK PATH

var add_smb_access_playbook string = playbook_dir + "add-smb-playbook.yml"
var remove_smb_access_playbook string = playbook_dir + "remove-smb-playbook.yml"
var install_smb_playbook string = playbook_dir + "install-smb-playbook.yml"
var uninstall_smb_playbook string = playbook_dir + "uninstall-smb-playbook.yml"

// ANSIBLE FUNCTION
// NFS FUNCTION

// ADD A NEW NFS ACCESS
/*
This function create a new nfs access, it require an list of ipaddress (string) to add in "/etc/exports" file
using a playbook called "add-nfs-playbook.yml", if ansible playbook cant be run it return an error else return nil value.
*/
func add_nfs_access(dests []string) error {

	ansiblePlaybookConnectionOptions := &ansible.AnsiblePlaybookConnectionOptions{}

	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{
		ExtraVars: map[string]interface{}{
			"dests": dests,
		},
	}

	ansible := &ansible.AnsiblePlaybookCmd{
		Playbook:          add_nfs_access_playbook,
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
	}
	if err := ansible.Run(); err != nil {
		return ERR_cant_run_add_nfs_access_playbook
	}
	return nil
}

// ADD A NEW SMB ACCESS
/*
This function create a new smb access, it require an list of users (see smb_user struct) to create it, set their password and add they to smb_group.

For do this this function run a playbook called "add-smb-playbook.yml", if ansible playbook cant be run it return an error else return nil value.
*/
func add_smb_access(users []smb_user) error {

	ansiblePlaybookConnectionOptions := &ansible.AnsiblePlaybookConnectionOptions{}

	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{
		ExtraVars: map[string]interface{}{
			"users": users,
		},
	}

	ansible := &ansible.AnsiblePlaybookCmd{
		Playbook:          add_smb_access_playbook,
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
	}
	if err := ansible.Run(); err != nil {
		return ERR_cant_run_add_smb_access_playbook
	}
	return nil
}

// REMOVE AN NFS ACCESS
/*
This function remove an nfs access, it require an list of ipaddress (string) to remove from "/etc/exports" file
using a playbook called "remove-nfs-playbook.yml", if ansible playbook cant be run it return an error else return nil value.
*/
func remove_nfs_access(dests []string) error {

	ansiblePlaybookConnectionOptions := &ansible.AnsiblePlaybookConnectionOptions{}

	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{
		ExtraVars: map[string]interface{}{
			"dests": dests,
		},
	}

	ansible := &ansible.AnsiblePlaybookCmd{
		Playbook:          remove_nfs_access_playbook,
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
	}
	if err := ansible.Run(); err != nil {
		return ERR_cant_run_remove_nfs_access_playbook
	}
	return nil
}

// REMOVE AN SMB ACCESS
/*
This function remove an smb access, it require an list of users (see smb_user struct) to remove.

For do this, this function run a playbook called "remove-smb-playbook.yml", if ansible playbook cant be run it return an error else return nil value.
*/
func remove_smb_access(users []string) error {

	ansiblePlaybookConnectionOptions := &ansible.AnsiblePlaybookConnectionOptions{}

	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{
		ExtraVars: map[string]interface{}{
			"users": users,
		},
	}

	ansible := &ansible.AnsiblePlaybookCmd{
		Playbook:          remove_smb_access_playbook,
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
	}
	if err := ansible.Run(); err != nil {
		return ERR_cant_run_remove_smb_access_playbook
	}
	return nil
}

//INSTALL NFS
/*
This function install all required nfs package and configure nfs env.

For do this is required and ip address src used for mount by nfs a volume to rexpose.

Actions:

- 1 install nfs

- 2 create exports directory

- 3 mount nfs volume by source in exports directory

- 4 create if it not exist /etc/exports

- 5 start nfs services

*/
func install_nfs_server(src string) error {

	ansiblePlaybookConnectionOptions := &ansible.AnsiblePlaybookConnectionOptions{}

	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{
		ExtraVars: map[string]interface{}{
			"src": src,
		},
	}

	ansible := &ansible.AnsiblePlaybookCmd{
		Playbook:          install_nfs_playbook,
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
	}
	if err := ansible.Run(); err != nil {
		return ERR_cant_run_install_nfs_playbook
	}
	return nil
}

//INSTALL SMB
/*
This function install all required smb package and configure smb env.

For do this is required and ip address src used for mount by nfs a volume to rexpose and the server workgroup name.

Actions:

- 1 install smb and nfs-utils

- 2 create sharing directory

- 3 mount nfs volume by source in sharing directory

- 4 create smb_group (GID=1005)

- 5 create if it not exist /etc/samba/smb.conf and configure it

- 6 start smb services

*/
func install_smb_server(src string, workgroup string) error {

	ansiblePlaybookConnectionOptions := &ansible.AnsiblePlaybookConnectionOptions{}

	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{
		ExtraVars: map[string]interface{}{
			"src":       src,
			"workgroup": workgroup,
		},
	}

	ansible := &ansible.AnsiblePlaybookCmd{
		Playbook:          install_smb_playbook,
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
	}
	if err := ansible.Run(); err != nil {
		return ERR_cant_run_install_smb_playbook
	}
	return nil
}

//UNINSTALL NFS
/*
This function uninstall all nfs package and delete nfs configurations.

Actions:

- 1 stop nfs services

- 2 uninstall nfs package

- 3 umount src nfs volume

- 4 remove exports directory

- 5 remove /etc/exports file

*/
func uninstall_nfs_server() error {

	ansiblePlaybookConnectionOptions := &ansible.AnsiblePlaybookConnectionOptions{}

	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{}

	ansible := &ansible.AnsiblePlaybookCmd{
		Playbook:          uninstall_nfs_playbook,
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
	}
	if err := ansible.Run(); err != nil {
		return ERR_cant_run_uninstall_nfs_playbook
	}
	return nil
}

//UNINSTALL SMB
/*
This function uninstall all nfs package and delete nfs configurations.

Actions:

- 1 stop smb services

- 2 get all users in smb_group (GID=1005)

- 3 remove users in smb_group (GID=1005)

- 4 remove smb_group

- 5 remove samba packages

- 6 umount src nfs volume

- 7 remove sharing directory

- 8 remove samba configuration directory /etc/samba

*/

func uninstall_smb_server() error {

	ansiblePlaybookConnectionOptions := &ansible.AnsiblePlaybookConnectionOptions{}

	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{}

	ansible := &ansible.AnsiblePlaybookCmd{
		Playbook:          uninstall_smb_playbook,
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
	}
	if err := ansible.Run(); err != nil {
		return ERR_cant_run_uninstall_smb_playbook
	}
	return nil
}

// API FUNCTIONS

/*
POST REQUEST FOR INSTALL NFS SERVER

- Src   string --> src ip address from which mount nfs volume

- Dests []string --> a list of ip address to which rexpose volume by nfs
*/
type post_server_nfs_request struct {
	Src   string   `json:"src"`
	Dests []string `json:"dests"`
}

/*
POST REQUEST FOR INSTALL SMB SERVER

- Src   string -->	src ip address from which mount nfs volume

- Workgroup string -->	name of smb server workgroup

- Users []smb_user --> a list of users see smb_user who can access by smb to rexposed volume
*/
type post_server_smb_request struct {
	Src       string     `json:"src"`
	Workgroup string     `json:"workgroup"`
	Users     []smb_user `json:"users"`
}

/*
SMB USER STRUCT

- User   string -->	username of smb user

- Pass string --> password of smb user
*/
type smb_user struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

/*
POST REQUEST FOR ADD SMB ACCESS

- Users []smb_user --> a list of users see smb_user who can access by smb to rexposed volume
*/
type server_smb_add_access_request struct {
	Users []smb_user `json:"users"`
}

/*
POST REQUEST FOR REMOVE SMB ACCESS

- Users []string --> a list of users username to be deleted by smb access
*/
type server_smb_remove_access_request struct {
	Users []string `json:"users"`
}

/*
POST REQUEST FOR ADD/REMOVE NFS ACCESS

- Dests []string --> a list of ip to be adds/removed by nfs access
*/
type server_nfs_access_request struct {
	Dests []string `json:"dests"`
}

/*
SIMPLE JSON RESPONSE

- Messagge string --> a simple text response messagge
- Status Status --> return operation status see Status
*/
type Response struct {
	Message string `json:"message"`
	Status  Status `json:"status"`
}

/*
SIMPLE Status tracker

- Situation string --> describe operation status
- Done bool --> if operation as be done without error
*/
type Status struct {
	Situation string `json:"situation"`
	Done      bool   `json:"done"`
}

/*
Avaible Status list

0:

Situation: "undone",

Done:      false,

1:
Situation: "done",

Done:      true,
*/
var status = []Status{
	{
		Situation: "undone",
		Done:      false,
	},
	{
		Situation: "done",
		Done:      true,
	},
}

/*
API FUNCTION INSTALL NFS REQUEST

This function require a json file in post_server_nfs_request format, see above.

It install an nfs server and configure first access next return operations status.

- 1 install_nfs_server

- 2 add_nfs_access
*/
func api_install_server_nfs(context *gin.Context) {
	var post_request post_server_nfs_request
	if err := context.BindJSON(&post_request); err != nil {
		context.IndentedJSON(http.StatusBadRequest, ERR_bad_json_format)
		return
	}
	if err := install_nfs_server(post_request.Src); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error(), Status: status[0]})
		return
	}
	if err := add_nfs_access(post_request.Dests); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error(), Status: status[0]})
		return
	}
	context.IndentedJSON(http.StatusOK, Response{Message: RES_installed, Status: status[1]})
}

/*
API FUNCTION INSTALL SMB REQUEST

This function require a json file in post_server_smb_request format, see above.

It install an smb server and configure first access next return operations status.

- 1 install_smb_server

- 2 add_smb_access
*/
func api_install_server_smb(context *gin.Context) {
	var post_request post_server_smb_request
	if err := context.BindJSON(&post_request); err != nil {
		context.IndentedJSON(http.StatusBadRequest, ERR_bad_json_format)
		return
	}
	if err := install_smb_server(post_request.Src, post_request.Workgroup); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error(), Status: status[0]})
		return
	}
	if err := add_smb_access(post_request.Users); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error(), Status: status[0]})
		return
	}
	context.IndentedJSON(http.StatusOK, Response{Message: RES_installed, Status: status[1]})
}

/*
API FUNCTION ADD NFS ACCESS REQUEST

This function require a json file in server_nfs_access_request format, see above.

It add new nfs access.

- 1 add_nfs_access
*/
func api_post_server_nfs_access(context *gin.Context) {
	var post_request server_nfs_access_request
	if err := context.BindJSON(&post_request); err != nil {
		context.IndentedJSON(http.StatusBadRequest, ERR_bad_json_format)
		return
	}

	if err := add_nfs_access(post_request.Dests); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error(), Status: status[0]})
		return
	}
	context.IndentedJSON(http.StatusOK, Response{Message: RES_added, Status: status[1]})
}

/*
API FUNCTION ADD SMB ACCESS REQUEST

This function require a json file in server_smb_add_access_request format, see above.

It add new smb access.

- 1 add_smb_access
*/
func api_post_server_smb_access(context *gin.Context) {
	var post_request server_smb_add_access_request
	if err := context.BindJSON(&post_request); err != nil {
		context.IndentedJSON(http.StatusBadRequest, ERR_bad_json_format)
		return
	}

	if err := add_smb_access(post_request.Users); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error(), Status: status[0]})
		return
	}
	context.IndentedJSON(http.StatusOK, Response{Message: RES_added, Status: status[1]})
}

/*
API FUNCTION REMOVE NFS ACCESS REQUEST

This function require a json file in server_nfs_access_request format, see above.

It remove old nfs access.

- 1 remove_nfs_access
*/
func api_patch_server_nfs_access(context *gin.Context) {
	var post_request server_nfs_access_request
	if err := context.BindJSON(&post_request); err != nil {
		context.IndentedJSON(http.StatusBadRequest, ERR_bad_json_format)
		return
	}

	if err := remove_nfs_access(post_request.Dests); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error(), Status: status[0]})
		return
	}
	context.IndentedJSON(http.StatusOK, Response{Message: RES_removed, Status: status[1]})
}

/*
API FUNCTION REMOVE SMB ACCESS REQUEST

This function require a json file in server_smb_remove_access_request format, see above.

It remove old smb access.

- 1 remove_smb_access
*/
func api_patch_server_smb_access(context *gin.Context) {
	var post_request server_smb_remove_access_request
	if err := context.BindJSON(&post_request); err != nil {
		context.IndentedJSON(http.StatusBadRequest, ERR_bad_json_format)
		return
	}

	if err := remove_smb_access(post_request.Users); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error(), Status: status[0]})
		return
	}
	context.IndentedJSON(http.StatusOK, Response{Message: RES_removed, Status: status[1]})
}

/*
API FUNCTION UNINSTALL NFS REQUEST

It uninstall nfs from system.

- 1 uninstall_nfs_server
*/
func api_uninstall_server_nfs(context *gin.Context) {
	if err := uninstall_nfs_server(); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error(), Status: status[0]})
		return
	}
	context.IndentedJSON(http.StatusOK, Response{Message: RES_uninstalled, Status: status[1]})
}

/*
API FUNCTION UNINSTALL SMB REQUEST

It uninstall smb from system.

- 1 uninstall_smb_server
*/
func api_uninstall_server_smb(context *gin.Context) {
	if err := uninstall_smb_server(); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error(), Status: status[0]})
		return
	}
	context.IndentedJSON(http.StatusOK, Response{Message: RES_uninstalled, Status: status[1]})
}

func main() {
	router := gin.Default()
	router.POST("/apiansible/server/nfs", api_install_server_nfs)
	router.POST("/apiansible/server/smb", api_install_server_smb)
	router.POST("/apiansible/server/nfs/access", api_post_server_nfs_access)
	router.POST("/apiansible/server/smb/access", api_post_server_smb_access)
	router.PATCH("/apiansible/server/nfs/access/", api_patch_server_nfs_access)
	router.PATCH("/apiansible/server/smb/access/", api_patch_server_smb_access)
	router.DELETE("/apiansible/server/nfs", api_uninstall_server_nfs)
	router.DELETE("/apiansible/server/smb", api_uninstall_server_smb)

	router.Run("0.0.0.0:4444")
}

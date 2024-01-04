package main
/*
NAME=APIANISBLE
AUTHOR=RH363
DATE=01/2024
COMPANY=SEEWEB
VERSION=1.0
*/
import (
	"errors"
	"net/http"

	ansible "github.com/febrianrendak/go-ansible"
	"github.com/gin-gonic/gin"
)

var ERR_bad_json_format error = errors.New("[ERROR] BAD JSON FORMAT")

var playbook_dir string = "/etc/ansible/playbooks/"

var add_nfs_access_playbook string = playbook_dir + "add-nfs-playbook.yml"
var remove_nfs_access_playbook string = playbook_dir + "remove-nfs-playbook.yml"
var install_nfs_playbook string = playbook_dir + "install-nfs-playbook.yml"
var uninstall_nfs_playbook string = playbook_dir + "uninstall-nfs-playbook.yml"

var add_smb_access_playbook string = playbook_dir + "add-smb-playbook.yml"
var remove_smb_access_playbook string = playbook_dir + "remove-smb-playbook.yml"
var install_smb_playbook string = playbook_dir + "install-smb-playbook.yml"
var uninstall_smb_playbook string = playbook_dir + "uninstall-smb-playbook.yml"

// START ANSIBLE FUNCTION
// START NFS FUNCTION
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
		return err
	}
	return nil
}

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
		return err
	}
	return nil
}

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
		return err
	}
	return nil
}

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
		return err
	}
	return nil
}

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
		return err
	}
	return nil
}

func install_smb_server(src string, workspace string) error {

	ansiblePlaybookConnectionOptions := &ansible.AnsiblePlaybookConnectionOptions{}

	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{
		ExtraVars: map[string]interface{}{
			"src":       src,
			"workspace": workspace,
		},
	}

	ansible := &ansible.AnsiblePlaybookCmd{
		Playbook:          install_smb_playbook,
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
	}
	if err := ansible.Run(); err != nil {
		return err
	}
	return nil
}

func uninstall_nfs_server() error {

	ansiblePlaybookConnectionOptions := &ansible.AnsiblePlaybookConnectionOptions{}

	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{}

	ansible := &ansible.AnsiblePlaybookCmd{
		Playbook:          uninstall_nfs_playbook,
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
	}
	if err := ansible.Run(); err != nil {
		return err
	}
	return nil
}

func uninstall_smb_server() error {

	ansiblePlaybookConnectionOptions := &ansible.AnsiblePlaybookConnectionOptions{}

	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{}

	ansible := &ansible.AnsiblePlaybookCmd{
		Playbook:          uninstall_smb_playbook,
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
	}
	if err := ansible.Run(); err != nil {
		return err
	}
	return nil
}

// END NFS FUNCTION
// START API FUNCTION

type post_server_nfs_request struct {
	Src   string   `json:"src"`
	Dests []string `json:"dests"`
}

type post_server_smb_request struct {
	Src       string     `json:"src"`
	Workspace string     `json:"workspace"`
	Users     []smb_user `json:"users"`
}

type smb_user struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type server_smb_add_access_request struct {
	Users []smb_user `json:"users"`
}

type server_smb_remove_access_request struct {
	Users []string `json:"users"`
}

type server_nfs_access_request struct {
	Dests []string `json:"dests"`
}

type Response struct {
	Message string `json:"message"`
}

func api_install_server_nfs(context *gin.Context) {
	var post_request post_server_nfs_request
	if err := context.BindJSON(post_request); err != nil {
		context.IndentedJSON(http.StatusBadRequest, ERR_bad_json_format)
		return
	}
	if err := install_nfs_server(post_request.Src); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}
	if err := add_nfs_access(post_request.Dests); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, post_request)
}

func api_install_server_smb(context *gin.Context) {
	var post_request post_server_smb_request
	if err := context.BindJSON(post_request); err != nil {
		context.IndentedJSON(http.StatusBadRequest, ERR_bad_json_format)
		return
	}
	if err := install_smb_server(post_request.Src, post_request.Workspace); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}
	if err := add_smb_access(post_request.Users); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, post_request)
}

func api_post_server_nfs_access(context *gin.Context) {
	var post_request server_nfs_access_request
	if err := context.BindJSON(post_request); err != nil {
		context.IndentedJSON(http.StatusBadRequest, ERR_bad_json_format)
		return
	}

	if err := add_nfs_access(post_request.Dests); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, post_request)
}

func api_post_server_smb_access(context *gin.Context) {
	var post_request server_smb_add_access_request
	if err := context.BindJSON(post_request); err != nil {
		context.IndentedJSON(http.StatusBadRequest, ERR_bad_json_format)
		return
	}

	if err := add_smb_access(post_request.Users); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, post_request)
}

func api_patch_server_nfs_access(context *gin.Context) {
	var post_request server_nfs_access_request
	if err := context.BindJSON(post_request); err != nil {
		context.IndentedJSON(http.StatusBadRequest, ERR_bad_json_format)
		return
	}

	if err := remove_nfs_access(post_request.Dests); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, post_request)
}

func api_patch_server_smb_access(context *gin.Context) {
	var post_request server_smb_remove_access_request
	if err := context.BindJSON(post_request); err != nil {
		context.IndentedJSON(http.StatusBadRequest, ERR_bad_json_format)
		return
	}

	if err := remove_smb_access(post_request.Users); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, post_request)
}

func api_uninstall_server_nfs(context *gin.Context) {
	if err := uninstall_nfs_server(); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, Response{Message: "uninstalled"})
}

func api_uninstall_server_smb(context *gin.Context) {
	if err := uninstall_smb_server(); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, Response{Message: "uninstalled"})
}

// END API FUNCTION
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

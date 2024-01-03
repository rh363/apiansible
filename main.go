package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	ansible "github.com/febrianrendak/go-ansible"
	"github.com/gin-gonic/gin"
)

var ERR_inventory_compromised error = errors.New("[PANIC] INVENTORY COMPROMISED TRYING RESTORE BY BACKUP")
var ERR_inventory_not_updated error = errors.New("[ERROR] INVENTORY NOT UPDATED, OLD INVENTORY RESTORED")
var ERR_inventory_not_restored error = errors.New("[PANIC] INVENTORY NOT UPDATED, CAN'T RESTORE OLD INVENTORY PLEASE RESTORE IT FROM \"" + inventory_backup_path + "\"")

var ERR_hosts_compromised error = errors.New("[PANIC] HOSTS COMPROMISED TRYING RESTORE BY BACKUP")
var ERR_hosts_not_updated error = errors.New("[ERROR] HOSTS NOT UPDATED, OLD INVENTORY RESTORED")
var ERR_hosts_not_restored error = errors.New("[PANIC] HOSTS NOT UPDATED, CAN'T RESTORE OLD INVENTORY PLEASE RESTORE IT FROM \"" + inventory_backup_path + "\"")

var ERR_server_not_found error = errors.New("[ERROR] SERVER NOT FOUND")
var ERR_server_already_exist error = errors.New("[ERROR] SERVER ALREADY EXISTS")

var ERR_bas_json_format error = errors.New("[ERROR] BAD JSON FORMAT")
var ansible_user string = "ansible-user"

var playbook_dir string = "/etc/ansible/playbooks/"

var add_nfs_access_playbook string = playbook_dir + "add-nfs-playbook.yml"
var remove_nfs_access_playbook string = playbook_dir + "remove-nfs-playbook.yml"
var install_nfs_playbook string = playbook_dir + "install-nfs-playbook.yml"
var uninstall_nfs_playbook string = playbook_dir + "uninstall-nfs-playbook.yml"

var hosts_file string = "hosts"
var inventory_file string = "hosts"
var inventory_backup_file string = "backup"
var hosts_backup_file string = "hosts.backup"

var hosts_path string = "/etc/" + hosts_file
var inventory_path string = "/root/" + inventory_file
var inventory_backup_path = "/etc/ansible/" + inventory_backup_file
var hosts_backup_path = "/etc/" + hosts_backup_file

// START FILE MANAGE FUNCTION
func ReadFile(path string) ([]string, error) { // READ A FILE CONTENT, REQUIRE FILE PATH, RETURN THE CONTENT AND AN ERROR

	File, err := os.Open(path)
	defer File.Close()
	if err != nil {
		return nil, err
	}

	reader := bufio.NewScanner(File)
	reader.Split(bufio.ScanLines)
	var txt []string

	for reader.Scan() {
		if reader.Text() != "" && reader.Text() != " " {
			//fmt.Println(reader.Text())
			txt = append(txt, reader.Text())
		}
	}

	return txt, nil
}

func WriteFile(path string, txt []string) error { // WRITE A NEW FILE IF IT DOESN'T EXIST, ELSE CREATE A NEW FILE, REQUIRE FILE PATH AND CONTENT, RETURN AN ERROR

	File, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	defer File.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	writer := bufio.NewWriter(File)

	for _, data := range txt {
		_, err = writer.WriteString(data + "\n")
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	writer.Flush()
	return nil
}

func insert_string(slice []string, index int, s string) []string {
	var ns []string
	for i, line := range slice {
		ns = append(ns, line)
		if i == index {
			ns = append(ns, s)
		}
	}
	return ns
}

func remove_string(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...) // it take content in slice until the index passed and append every content after the index to it, ... is used for pass element in second slice one by one
}

// END FILE MANAGE FUNCTION
// START ANSIBLE FUNCTION
func add_server(namespace string, ip string) error {
	inventory_txt, err := ReadFile(inventory_path)
	if err != nil {
		return err
	}
	for _, line := range inventory_txt {
		if line == namespace+"-"+ip+" ansible_user=ansible-user" {
			return ERR_server_already_exist
		}
	}
	hosts_txt, err := ReadFile(hosts_path)
	if err != nil {
		return err
	}

	if err := WriteFile(inventory_backup_path, inventory_txt); err != nil {
		return err
	}

	hostname_registered := false
	for _, host := range hosts_txt {
		if host == ip+" "+namespace+"-"+ip {
			hostname_registered = true
		}
	}

	if !hostname_registered {
		if err := WriteFile(hosts_path, []string{ip + " " + namespace + "-" + ip}); err != nil {
			return err
		}
	}

	var new_inventory_txt []string
	for i, line := range inventory_txt {
		if line == "["+namespace+"]" {

			new_inventory_txt = insert_string(inventory_txt, i, namespace+"-"+ip+" ansible_user=ansible-user")

			if err := os.Remove(inventory_path); err != nil {
				return err
			}

			for i := 0; i < 10; i++ {
				if err := WriteFile(inventory_path, new_inventory_txt); err == nil {
					return nil
				}
				fmt.Println(ERR_inventory_compromised)
				for i := 0; i < 10; i++ {
					if err := WriteFile(inventory_path, inventory_txt); err == nil {
						return ERR_inventory_not_updated
					}
				}
				return ERR_inventory_not_restored
			}
		}
	}
	if err := WriteFile(inventory_path, []string{
		"[" + namespace + "]",
		namespace + "-" + ip + " ansible_user=ansible-user",
	}); err != nil {
		return ERR_inventory_not_updated
	}

	return nil
}

func remove_server(namespace string, ip string) error {
	inventory_txt, err := ReadFile(inventory_path)
	if err != nil {
		return err
	}
	if err := WriteFile(inventory_backup_path, inventory_txt); err != nil {
		return err
	}

	hosts_txt, err := ReadFile(hosts_path)
	if err != nil {
		return err
	}
	if err := WriteFile(hosts_backup_path, hosts_txt); err != nil {
		return err
	}

	var removed_server_txt []string
	var removed_host_txt []string
	var removed_server_group_txt []string

	for i, host := range hosts_txt {
		if host == ip+" "+namespace+"-"+ip {
			removed_host_txt = remove_string(hosts_txt, i)
			if err := os.Remove(hosts_path); err != nil {
				return err
			}

			for i := 0; i < 10; i++ {
				if err := WriteFile(hosts_path, removed_host_txt); err != nil {
					fmt.Println(ERR_hosts_compromised)
					for i := 0; i < 10; i++ {
						if err := WriteFile(hosts_path, hosts_txt); err == nil {
							return ERR_hosts_not_updated
						}
					}
					return ERR_hosts_not_restored
				}
				break
			}
		}
	}

	for i, line := range inventory_txt {
		if line == namespace+"-"+ip+" ansible_user=ansible-user" {

			removed_server_txt = remove_string(inventory_txt, i)

			for _, line := range removed_server_txt {
				if line != "["+namespace+"]" && strings.Contains(line, namespace) {

					if err := os.Remove(inventory_path); err != nil {
						return err
					}
					for i := 0; i < 10; i++ {
						if err := WriteFile(inventory_path, removed_server_txt); err == nil {
							return nil
						}
						fmt.Println(ERR_inventory_compromised)
						for i := 0; i < 10; i++ {
							if err := WriteFile(inventory_path, inventory_txt); err == nil {
								return ERR_inventory_not_updated
							}
						}
						return ERR_inventory_not_restored
					}

				}
			}
			for i, line := range removed_server_txt {
				if line == "["+namespace+"]" {

					removed_server_group_txt = remove_string(removed_server_txt, i)

					if err := os.Remove(inventory_path); err != nil {
						return err
					}
					for i := 0; i < 10; i++ {
						if err := WriteFile(inventory_path, removed_server_group_txt); err == nil {
							return nil
						}
						fmt.Println(ERR_inventory_compromised)
						for i := 0; i < 10; i++ {
							if err := WriteFile(inventory_path, inventory_txt); err == nil {
								return ERR_inventory_not_updated
							}
						}
						return ERR_inventory_not_restored
					}

				}
			}
		}
	}
	return ERR_server_not_found
}

// START NFS FUNCTION
func add_nfs_access(target string, dests []string) error {

	ansiblePlaybookConnectionOptions := &ansible.AnsiblePlaybookConnectionOptions{
		User: ansible_user,
	}

	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{
		ExtraVars: map[string]interface{}{
			"target": target,
			"dests":  dests,
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

func remove_nfs_access(target string, dests []string) error {

	ansiblePlaybookConnectionOptions := &ansible.AnsiblePlaybookConnectionOptions{
		User: ansible_user,
	}

	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{
		ExtraVars: map[string]interface{}{
			"target": target,
			"dests":  dests,
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

func install_nfs_server(target string, src string) error {

	ansiblePlaybookConnectionOptions := &ansible.AnsiblePlaybookConnectionOptions{
		User: ansible_user,
	}

	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{
		ExtraVars: map[string]interface{}{
			"target": target,
			"src":    src,
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

func uninstall_nfs_server(target string) error {

	ansiblePlaybookConnectionOptions := &ansible.AnsiblePlaybookConnectionOptions{
		User: ansible_user,
	}

	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{
		ExtraVars: map[string]interface{}{
			"target": target,
		},
	}

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

// END NFS FUNCTION
// START API FUNCTION
type target struct {
	Namespace string `json:"namespace"`
	Ip        string `json:"ip"`
}

type serverlist struct {
	Namespace string   `json:"namespace"`
	Ips       []string `json:"ips"`
}

type post_server_nfs_request struct {
	Target target   `json:"target"`
	Src    string   `json:"src"`
	Dests  []string `json:"dests"`
}

type server_nfs_access_request struct {
	Target target   `json:"target"`
	Dests  []string `json:"dests"`
}

type Response struct {
	Message string `json:"message"`
}

func get_server(context *gin.Context) {
	var server_list []serverlist
	var namespace string
	var ips []string
	inventory, err := ReadFile(inventory_path)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, Response{Message: err.Error()})
		return
	}
	for _, line := range inventory {
		fmt.Println(line)
		if strings.Contains(line, "[") {
			server_list = append(server_list, serverlist{Namespace: namespace, Ips: ips})
			ips = nil
			namespace = line
		} else {
			ips = append(ips, strings.Split(line, " ")[0])
		}
	}

	server_list = append(server_list, serverlist{Namespace: namespace, Ips: ips})

	context.IndentedJSON(http.StatusOK, server_list[1:])
}

func post_server(context *gin.Context) {
	var post_request serverlist
	if err := context.BindJSON(&post_request); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: ERR_bas_json_format.Error()})
		return
	}

	for _, ip := range post_request.Ips {
		if err := add_server(post_request.Namespace, ip); err != nil {
			context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error()})
		}
	}
	context.IndentedJSON(http.StatusOK, post_request)
}

func post_server_nfs(context *gin.Context) {
	var post_request post_server_nfs_request
	if err := context.BindJSON(post_request); err != nil {
		context.IndentedJSON(http.StatusBadRequest, ERR_bas_json_format)
		return
	}
	if err := install_nfs_server(post_request.Target.Namespace+"-"+post_request.Target.Ip, post_request.Src); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}
	if err := add_nfs_access(post_request.Target.Namespace+"-"+post_request.Target.Ip, post_request.Dests); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, post_request)
}

func post_server_nfs_access(context *gin.Context) {
	var post_request server_nfs_access_request
	if err := context.BindJSON(post_request); err != nil {
		context.IndentedJSON(http.StatusBadRequest, ERR_bas_json_format)
		return
	}

	if err := add_nfs_access(post_request.Target.Namespace+"-"+post_request.Target.Ip, post_request.Dests); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, post_request)
}

func patch_server_nfs_access(context *gin.Context) {
	var post_request server_nfs_access_request
	if err := context.BindJSON(post_request); err != nil {
		context.IndentedJSON(http.StatusBadRequest, ERR_bas_json_format)
		return
	}

	if err := remove_nfs_access(post_request.Target.Namespace+"-"+post_request.Target.Ip, post_request.Dests); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, post_request)
}

func delete_server_nfs(context *gin.Context) {
	namespace := context.Param("namespace")
	ip := context.Param("ip")

	if err := uninstall_nfs_server(namespace + "-" + ip); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, Response{Message: namespace + "-" + ip})
}

func delete_server(context *gin.Context) {
	namespace := context.Param("namespace")
	ip := context.Param("ip")

	if err := remove_server(namespace, ip); err != nil {
		context.IndentedJSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}

	context.IndentedJSON(http.StatusOK, Response{Message: namespace + "-" + ip})
}

// END API FUNCTION
func main() {
	router := gin.Default()
	router.GET("/apiansible/inventory/server", get_server)
	router.POST("/apiansible/inventory/server/", post_server)
	router.POST("/apiansible/inventory/server/nfs", post_server_nfs)
	router.POST("/apiansible/inventory/server/nfs/access", post_server_nfs_access)
	router.PATCH("/apiansible/inventory/server/nfs/access/", patch_server_nfs_access)
	router.DELETE("/apiansible/inventory/server/:namespace/:ip/nfs/", delete_server_nfs)
	router.DELETE("/apiansible/inventory/server/:namespace/:ip", delete_server)

	router.Run("0.0.0.0:6666")
}

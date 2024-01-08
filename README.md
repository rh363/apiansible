
# APIANSIBLE

This project is developed for rexpose an nfs volume from a secure network to an insecure netwok. \
There are two way to rexpose this volume: 
- 1 NFS
- 2 SMB

## How it work?
This project for work use some ansible playbook who can be runned by an go api called "apiansible". \
Using our "recipes" is possible install an nfs/smb server modify the access to this server and if is required delete all pacakge, configuration files and volume data directory with a simple http api request.


## Installation

for install apiasible you have to claim this steps:

- 1 install ansible
- 2 configure ansible
- 3 pull ansible playbook for configure your server

Install this project is very easy and fast with this repo. \
You only have to mannualy download ansible using the [official ansible doc](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html). \
I strongly reccomend you to install it on rhel like os for example almalinux or centos. \
Next you only have to create or change the following file to the ones provided in this repo:
- /etc/ansible/ansible.cfg
- /etc/ansible/hosts
Now you can finnaly pull installation playbook:

```bash
  ansible-pull -U https://github.com/rh363/apiansible.git
```
Can you see all operations performed by this pull in "local.yml"
Now you have the api running on service apiansible.service and listenning on port 4444.
## Requests
The http request implemented are separated in nfs and smb requests:
## NFS Request:

The following requests are used for install and manage an nfs server. \
The server mount an nfs volume and rexpose it via nfs.

**INSTALL NFS SERVER:**

for install nfs client this POST request is required.

|POST Request | **Request body** | **Explain** |
| --- | --- | --- |
| \*server ip address\*:4444/apiansible/server/nfs | {  <br>"src":"\*nfs volume ip to rexpose\*",  <br>"dests":\[  <br>"\*ip who can access to nfs volume\*"  <br>\]  <br>} | this request required two argoument:  <br>src is the nfs ip address who our machine mount.  <br>dests is a list of ip address who can access to the rexposed nfs volume |

**ADD NFS ACCESS:**

this POST request add another access to our nfs volume.

| POST Request | Request body | Explain |
| --- | --- | --- |
| \*server ip address\*:4444/apiansible/server/nfs/access | {  <br>"dests":\[  <br>"\*new access\*"  <br>\]  <br>} | this request require only a list of ip address to add in /etc/exports file.  <br>this list is called dests. |

**REMOVE NFS ACCESS:**

This PATCH request remove an nfs access to our rexposed volume.

|PATCH Request | **Request body** | Explain |
| --- | --- | --- |
| \*server ip address\*:4444/apiansible/server/nfs/access/ | {  <br>"dests": \[  <br>"\*access to remove\*"  <br>\]  <br>} | this function require a list of ip to remove from /etc/exports.  <br>this list is called dests. |

**UNINSTALL NFS SERVER:**

use this DELETE request for uninstall nfs server package and configuration from system.

|DELETE Request | Explain |
| --- | --- |
| \*server ip address\*/apiansible/server/nfs | dont require any body.  <br>use this request for clean server from nfs package and configuration. |

## SMB Requests:

The following requests are used for install and manage an smb server. \
The server mount an nfs volume and rexpose it via smb.

**INSTALL SMB SERVER:**

This POST request is used for install an smb server.

for use this request you must know our user type struct:

``` json
"user" : {
    "user" : "username","pass" : "password"
}

 ```

| POST Request | **Request body** | Explain |
| --- | --- | --- |
| \*server ip address\*:4444/apiansible/server/smb | {  <br>"src": "nfs volume ip to rexpose",  <br>"workgroup": "\*workgroup name\*",  <br>"users":\[  <br>{"user": "\*username\*","pass":"\*password\*"}  <br>\]  <br>} | this request require some argument:  <br>the nfs ip address src.  <br>the smb workgroup name.  <br>a list of user who can access to the smb volume described using the above format.  <br>this list is called users. |

**ADD SMB ACCESS:**

This POST Request is used to add some access to our smb server.\
for use this request you must know our user type struct:

``` json
"user" : {
    "user" : "username","pass" : "password"
}

 ```

|POST Request | **Request body** | Explain |
| --- | --- | --- |
| \*server ip address\*:4444/apiansible/server/smb/access | {  <br>"users":\[  <br>{"user": "\*username\*","pass":"\*password\*"}  <br>\]  <br>} | this request require some argument:  <br>a list of user who can access to the smb volume described using the above format.  <br>this list is called users. |

**REMOVE SMB ACCESS:**

This PATCH request is used for remove and nfs access from the smb volume.

|PATCH Request | **Request body** | Explain |
| --- | --- | --- |
| \*server ip address\*/apiansible/server/smb/access/ | {  <br>"users":\[  <br>"\*username\*"  <br>\]  <br>} | this request require some argument:  <br>a list of user who can access to the smb volume, for define the user to remove only username is required.  <br>this list is called users. |

**UNINSTALL SMB SERVER:**

use this DELETE request for uninstall smb server package and configuration from system.

|DELETE Request | Explain |
| --- | --- |
| \*server ip address\*:4444/apiansible/server/smb/ | dont require any body.  <br>use this request for clean server from smb package and configuration. |


## Used Projects

- For api build: [gin-gonic](https://github.com/gin-gonic/gin)
- Ansible project: [Ansible](https://github.com/ansible/ansible)
## Authors

- [Massaroni alex](https://www.github.com/rh363)
- [Vona Daniele](https://github.com/danielv99)

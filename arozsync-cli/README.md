# arozsync-cli

This command line interface tool is designed to sync arozos files from WebDAV interface to local disk and vice-versa.

## Build
```
cd arozsync-cli
go build
```

## Usage
```
  -c string
        The configuration file path (default "config.json")
  -clean
        [DANGER] Execute system cleaning to remove deleted file backups
  -ld
        Enable remove local file by deleteing in remote file system (default true)
  -lw
        Enable remote changed to overwrite local changes (default true)
  -pass string
        Password for authentication
  -rd
        Enable remove remote file by deleting in local file system
  -rw
        Enable upload local change to remote file system. Set this to false for DOWNLOAD FROM REMOTE ONLY. (default true)
  -user string
        Username for authentication
```

For the most basic sync setup, you will need a config.json file like this.
```
{
 "ServerAddr": "127.0.0.1",
 "Port": 8080,
 "SyncInterval": 300,
 "UseHTTPs": false,
 "Folders": [
  {
   "RemoteRootID": "user:/",
   "RemoteFolder": "Desktop/",
   "LocalFolder": "./tmp"
  }
 ]
}
```

To generate this file locally, you can use ```./arozsync-cli``` with no paramters. A new config.json file will be generated.

Edit the json file to add your server information. Once it is done, start arozsync-cli with the following minimal paramters. The username and password is identical to the one you used to login to your ArozOS host.
```
./arozsync-cli -user="YOUR USER NAME" -pass="YOUR PASSWORD"
```
Then the tool will start to sync everything from your user:/Desktop folder to your ./tmp folder. 

### Remove Local file from Remote File System

This option is enabled by default. When a synchronized file is being deleted on ArozOS via Web UI, the local file will be moved to a ```.trash``` folder local at the same parent folder of the deleted file. These files will not be automatically deleted or uploaded to the server side again. If you want to restore them, you need to manually copy it from the .trash folder back to its original location.

If you want to clear all the backup files, simply use the following command

**⚠️ This operation cannot be undone**

```
./arozsync-cli -clean
```

### Remove Remote file from Local File System

**⚠️ Enabling this parameter is dangerous. Use with your own risk.**

If you want to allow users to remove remote file on ArozOS Server, you can set ```-rd=true``` in the startup parameter and this will allow user to remove remote file by deleting local file. All the deleted files are stored inside a hidden folder named ```.trash``` on the ArozOS server side. You can clear up the remote trash file using the Trash Bin WebApp in ArozOS


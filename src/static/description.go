package static

var Manual = `dejavuDB(1)		Manual pager		dejavuDB(1)

NAME
	dejavuDB console - an interface to the dejavuDB command

OPTIONS
    	dejavuDB [options] ...
    	dejavuDB -r      role          Specify role. Default as router, options: shard, client, full ...
       	dejavuDB -ip     listener_ip   Specify router ip. Default as stand alone router ...
       	dejavuDB -p      password      Specify password. Default as empty ...
       	dejavuDB -host   hostIP        Specify hosting port, Defult as localhost:8080 ...
        dejavuDB -disk                 Specify to save a copy to disk or not...
        dejavuDB -dr     path          Specify to read data from disk or not and the path...
        dejavuDB -sc     bits          specify to use securite connection and the bit width...

DESCRIPTION
        dejavuDB is a NoSql database written in golang. It has horizontal scalability and stores data in map[string]interface{} format.                                              

COMMANDS
        atop        The  program atop is an interactive monitor to view the load on a Linux system.  
                    It shows the occupation of the most critical hardware resources (from a perfor‐mance point of view) on system level, 
                    i.e. cpu, memory, disk and network.

        cat

        cp
	
        df

        dstat

        find

        free

        last

        mv

        netstat

        rm

        sort

        w           w  displays  information  about the users currently on the machine, and their processes.
                    The header shows, in this order, the current time, how long the system has been running,
                    how many users are currently logged on, and the system load averages.

        top        

        tar.xz      Compress Bytes in the .tar.xz format.

        tar         Compress Bytes in the .tar format.

        xz          Compress Bytes in the .xz format.
`

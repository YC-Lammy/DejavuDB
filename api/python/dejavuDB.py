import os, sys,io, socket,json, joblib
from tempfile import NamedTemporaryFile

def connect(host:str,port:str,username=None,password=None):
    conn = socket.create_connection((host,port))
    connfile = conn.makefile(mode='rw', encoding='utf-8')
    connfile.write("""{"role":"client","username":"{us}","password":"{pas}"}\n""".format(us=username,pas=password))
    config = connfile.readline() # config file from router
    return dejavuDB(conn,connfile,config)

class dejavuDB:
    def __init__(self,conn,connfile,config) -> None:
        self.conn = conn
        self.connfile = connfile
        self.config = config # config file from router
        self.SQL = DBSQL(conn,connfile)
        self.ML = DBML(conn, connfile)

    def close(self) -> None:
        self.conn.close()

    def Get(self,keys:str)-> str:
        "keys: syntax 'key0.key1.key2...'"

        self.connfile.write(f"Get {keys}"+0x00)
        return self.connfile.readline()

    def Set(self,keys:str,value,datatype:str):

        if datatype not in ["string","int","float","bool","bytes","[]byte"
        ,"[]string","[]int","[]float","[]bool","[][]byte"]:return ValueError("expected golang type, got "+datatype)

        self.connfile.write(f"Set {keys} {str(value)} {datatype}"+0x00)
        self.connfile.readline()

    def Clone(self,target_keys:str,destination_keys:str):
        self.connfile.write(f"Clone {target_keys} {destination_keys}\n")
        self.connfile.readline()

    def sendSQL(self,sql:str):
        self.connfile.write("SQL "+sql+0x00)
        return self.connfile.readline()

class DBSQL:
    def __init__(self,conn,connfile):
        self.conn = conn
        self.connfile = connfile
    
    def send(self, sql:str):
        self.connfile.write("SQL "+sql+0x00)
        return self.connfile.readline()
    
    def create_table(self,name:str,*argv:str):
        if len(argv) >1:
            for i in range(len(argv)):
                if i == 0:
                        continue
                if argv[i] in ["int","int8","int16","int32","varchar"]:
                    argv[i-1]+= " " +argv[i]
                    del argv[i]
            headers = argv.join(",")
        elif len(argv)== 1:
            headers = argv[0]
        else:
            raise Exception("atleast one arg needed")
        self.send("CREATE TABLE "+name+" ("+headers+")")

class DBML:
    def __init__(self,conn,connfile) -> None:
        self.conn = conn
        self.connfile = connfile
        self.keras_model = None
    
    def keras_compile(self,model,strategy = None,
    optimizer='rmsprop', loss=None, metrics=None, loss_weights=None,
    weighted_metrics=None, run_eagerly=None, steps_per_execution=None, **kwargs):
        model.fit # check if fit function exist
        if strategy != None:
            with strategy.scope():
                try:
                    model.compile(optimizer,loss,metrics,loss_weights,weighted_metrics,run_eagerly,steps_per_execution,kwargs)
                    self.keras_model = model
                except:
                    model = model()
                    model.compile(optimizer,loss,metrics,loss_weights,weighted_metrics,run_eagerly,steps_per_execution,kwargs)
                    self.keras_model = model
        else:
            try:
                    model.compile(optimizer,loss,metrics,loss_weights,weighted_metrics,run_eagerly,steps_per_execution,kwargs)
                    self.keras_model = model
            except:
                    model = model()
                    model.compile(optimizer,loss,metrics,loss_weights,weighted_metrics,run_eagerly,steps_per_execution,kwargs)
                    self.keras_model = model

    def tensorflow_compile(self,model):
        pass

    def save(self,name=None):
        with NamedTemporaryFile() as temp:

            if self.keras_model != None:
                self.keras_model.save(str(temp.name),save_format = "tf")
                modeltype= "keras"
            else:
                raise Exception("No model provide, please compile model using the compile function")

            json.dumps({"name":name,"type":modeltype,"model":temp.read()})

        


def formatData(data:str):

    if data[0] == "{" and data[-1] == "}":
        return json.loads(data)

    elif data[0]== "[":
        a = data[1:-1].split(",")
        try :
            return [int(i) if round(float(i)) == float(i) else float(i) for i in a]
        except: pass

        try: 
            return [bool(i) for i in a]
        except: pass

        return a

    else:
        try:return int(data) if round(float(data)) == float(data) else float(data)
        except: pass
        try:return bool(data)
        except:pass
        return data
    



        

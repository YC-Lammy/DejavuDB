import os, sys, socket,json

class dejavuDB:
    def __init__(self,host:str,port:int,username:str,password:str) -> None:

        self.conn = socket.create_connection((host,port))
        self.connfile = self.conn.makefile(mode='rw', encoding='utf-8')
        self.connfile.write("""{"role":"client","username":"{us}","password":"{pas}"}\n""".format(us=username,pas=password))
        self.config = self.connfile.readline() # config file from router

    def close(self) -> None:
        self.conn.close()

    def Get(self,keys:str)-> str:
        "keys: syntax 'key0.key1.key2...'"

        self.connfile.write(f"Get {keys}\n")
        return self.connfile.readline()

    def Set(self,keys:str,value,datatype:str):

        if datatype not in ["string","int","float","bool","bytes","[]byte"
        ,"[]string","[]int","[]float","[]bool","[][]byte"]:return ValueError("expected golang type, got "+datatype)

        self.connfile.write(f"Set {keys} {str(value)} {datatype}\n")
        self.connfile.readline()

    def Clone(self,target_keys:str,destination_keys:str):
        self.connfile.write(f"Clone {target_keys} {destination_keys}\n")
        self.connfile.readline()

    def sendSQL(self,sql:str):
        sql = sql.replace("\n"," ")
        self.connfile.write(sql+"\n")
        return self.connfile.readline()

    def formatData(self,data:str):

        if data[0] == "{":
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
            except:pass
            try:return bool(data)
            except:pass
            return data
            


        

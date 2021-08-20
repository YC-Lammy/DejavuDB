import os, sys, time
import socket, joblib
from tempfile import NamedTemporaryFile
from pathlib import Path
import multiprocessing
import threading

os.chdir(os.path.join(Path.home(),"dejavuDB","ML"))

header=64
method = 'UTF-8'
ender = "!END_OF_SERVICE"
WaitGroup = 0
resources = multiprocessing.Manager().dict() # map[string]multiprocessing.Lock
"""
if tf.config.list_physical_devices("TPU"):
    strategy = tf.distribute.TPUStrategy
elif tf.config.list_physical_devices("GPU"):
    strategy = tf.distribute.MirroredStrategy
else:
    strategy = tf.distribute.get_strategy()
"""

class pipline:
    def __init__(self,dict) -> None:

        self.model_name = dict["name"]

        self.ExampleGen = dict["ExampleGen"]
        self.StatisticsGen = dict["StatisticsGen"]
        self.SchemaGen = dict["SchemaGen"]
        self.ExampleValidator = dict["ExampleValidator"]
        self.Transform = dict["Transform"]
        self.Tuner = dict["Tuner"]
        self.Trainer = dict["Trainer"]
        self.Evaluator = dict["Evaluator"]
        self.InfraValidator = dict["InfraValidator"]
        self.Pusher = dict["Pusher"]
    
    def push(self,data):
        pass

    def save(self):
        joblib.dump(self,self.model_name+".pipline")

class SchemaGen:
    """schemaGen generates a schema containing features in dictionary

    feature {
        name: "age"
        value_count {
            min: 1
            max: 1
        }
        type: FLOAT
        presence {
            min_fraction: 1
            min_count: 1
        }
    }
    """
    def __init__(self,dict) -> None:
        self.features = dict["SchemaGen"]["features"]

class tenserflow_model:
    locked = False

    def __init__(self,model_path=None):
        self.model_path = model_path

def send(msg,conn):
    conn.send(msg.encode(method))

def recv(conn):
    msg_l = conn.recv(header).decode(method)
    if msg_l:
        msg_l = int(msg_l)
        msg = conn.recv(msg_l).decode(method)
    return msg


def main_handler(conn,msg,resources):
    import io, json, joblib
    import tensorflow as tf
    import numpy as np

    def send(msg,conn):
        conn.send(msg.encode(method))

    def sprint(*args, **kwargs):
        sio = io.StringIO()
        print(*args, **kwargs, file=sio)
        return sio.getvalue()

    def load_model(name:str):
        model = tf.keras.models.load_model(name)
        return model

    com_dict = json.loads(msg)
    com = com_dict["command"]
    if com == "CREATE_MODEL": #script to create and save model
        try:
            exec(msg[14:])
            send("sucess",conn)

        except Exception as e:
            send(str(e),conn)

    elif com == "CREATE_PIPLINE":
        pass

    elif com == "PREDICT": # PREDICT $model_name values
        if com_dict['name'] not in resources["model_names"]:
            send("model "+com_dict['name']+" does not exist", conn)
            return
        model = load_model(com_dict["name"])
        send(sprint(model.predict(com_dict["data"])),conn)

    elif com == "PUSH_PIPLINE":
        pass

    conn.close()

def spawn_process(conn,msg):
    global WaitGroup
    WaitGroup += 1
    try:
        parent_conn, child_conn = multiprocessing.Pipe()
        p = multiprocessing.Process(target=main_handler, args=(child_conn,msg,resources))
        p.start()
        a = parent_conn.recv()
        p.join()
        try:
            send(a, conn)
        except Exception as e:
            send(str(e),conn)
    except Exception as e:
        send(str(e),conn)

    WaitGroup = WaitGroup - 1

def main():
    soc =  socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    soc.bind(("localhost", "3247"))
    soc.listen()
    conn, addr = soc.accept() #since threading in pyhton does not actually uncrease speed i decided to use only one connection
    while True:
        msg = recv(conn) # msg in json format
        if msg == ender:
            while WaitGroup >0:
                time.sleep(0.25)
                continue
            break
        threading.Thread(target=spawn_process,args=(conn,msg)).start() #spawn new thread

if __name__ == "__main__":
    main()
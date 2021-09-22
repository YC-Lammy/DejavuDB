DB.Graph = class{
    constructor(path){
        this.IsGraph = true
        this.path = path
        if (path == undefined){
            this.path = call_go_fn("create","graph")
        }
    }
    addVertex(vertex){
        if (vertex.IsGraphVertex != true){
            throw "expected class GraphVertex";
        }
        call_go_fn("types")
    }
}

DB.GraphVertex = class{
    constructor(path, isSaved){
        this.IsGraphVertex = true
        this.values = {}
        this.path = path
        if (path == undefined){
            this.path = call_go_fn("create", "graphvertex")
        }
    }
    setValue(key, value){
        this.values[key] = value
        call_go_fn("value","callfn",this.path,"setvalue",key, value)
    }
    getValue(key){
        return call_go_fn("value","callfn",this.path,"getvalue",key)
    }
}
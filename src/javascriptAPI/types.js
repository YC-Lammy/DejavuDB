GoFunction = class{
    constructor(script, path){
        if (path == undefined){
            this.path = call_go_fn("create","Fn",script)
        } else{
            this.path = path
        }
        this.IsGoFunction = true;
    }
    toString(){
        return this.path
    }
    run(...any){
        return call_go_fn("value","call",this.path, ...any);
    }
}

GoString = class{
    constructor(str,path){
        if (typeof(str) != 'string'){
            throw "GoString constructor expected string"
        }
        if (path == undefined){
            this.path=  call_go_fn("create", "string", str)
        }
    }
    toString(){
        return this.path
    }
}


Graph = class{
    constructor(path){

        this.IsGraph = true
        this.path = path
        if (path == undefined){
            this.path = call_go_fn("create","graph")
        }
    }
    toString(){
        return this.path
    }
    addVertex(vertex){
        if (vertex.IsGraphVertex != true){
            throw "expected class GraphVertex";
        }
        call_go_fn("value","method",this.path,"AddVertex",vertex.toString())
    }
}

GraphVertex = class{
    constructor(path){

        this.IsGraphVertex = true
        this.values = {}
        this.path = path
        if (path == undefined){
            this.path = call_go_fn("create", "graphvertex")
        }
    }
    toString(){
        return this.path
    }
    setValue(key, value){
        this.values[key] = value
        call_go_fn("value","method",this.path,"SetField",key, value)
    }
    getValue(key){
        return call_go_fn("value","method",this.path,"GetField",key)
    }
}

GoInt = class{
    constructor(str,path){
        if (typeof(str) != 'number'|| path != undefined){
            throw "GoString constructor expected number"
        }
        if (path == undefined){
            this.path=  call_go_fn("create", "int", str)
        } else {
            this.path = path
        }

    }
    toString(){
        return this.path
    }
}

GoInt8 = class{
    constructor(str,path){
        if (typeof(str) != 'number'|| path != undefined){
            throw "GoString constructor expected number"
        }
        if (path == undefined){
            this.path=  call_go_fn("create", "int8", str)
        } else {
            this.path = path
        }
    }
    toString(){
        return this.path
    }
}

GoInt16 = class{
    constructor(str,path){
        if (typeof(str) != 'number'|| path != undefined){
            throw "GoString constructor expected number"
        }
        if (path == undefined){
            this.path=  call_go_fn("create", "int16", str)
        } else {
            this.path = path
        }
    }
    toString(){
        return this.path
    }
}

GoInt32 = class{
    constructor(str,path){
        if (typeof(str) != 'number'|| path != undefined){
            throw "GoString constructor expected number"
        }
        if (path == undefined){
            this.path=  call_go_fn("create", "int32", str)
        } else {
            this.path = path
        }
    }
    toString(){
        return this.path
    }
}

GoInt64 = class{
    constructor(str,path){
        if (typeof(str) != 'number'|| path != undefined){
            throw "GoString constructor expected number"
        }
        if (path == undefined){
            this.path=  call_go_fn("create", "int64", str)
        } else {
            this.path = path
        }
    }
    toString(){
        return this.path
    }
}

GoInt128 = class{
    constructor(str,path){
        if (typeof(str) != 'bigint'|| path != undefined){
            throw "GoString constructor expected bigint"
        }
        if (path == undefined){
            this.path=  call_go_fn("create", "int128", str.toString())
        } else {
            this.path = path
        }
    }
    toString(){
        return this.path
    }
}

GoUint = class{
    constructor(str,path){
        if (typeof(str) != 'number'|| path != undefined){
            throw "GoString constructor expected number"
        }
        if (path == undefined){
            this.path=  call_go_fn("create", "uint", str)
        } else {
            this.path = path
        }
    }
    toString(){
        return this.path
    }
}

GoUint8 = class{
    constructor(str,path){
        if (typeof(str) != 'number'|| path != undefined){
            throw "GoString constructor expected number"
        }
        if (path == undefined){
            this.path=  call_go_fn("create", "uint8", str)
        } else {
            this.path = path
        }
    }
    toString(){
        return this.path
    }
}

GoUint16 = class{
    constructor(str,path){
        if (typeof(str) != 'number'|| path != undefined){
            throw "GoString constructor expected number"
        }
        if (path == undefined){
            this.path=  call_go_fn("create", "uint16", str)
        } else {
            this.path = path
        }
    }
    toString(){
        return this.path
    }
}

GoUint32 = class{
    constructor(str,path){
        if (typeof(str) != 'number'|| path != undefined){
            throw "GoString constructor expected number"
        }
        if (path == undefined){
            this.path=  call_go_fn("create", "uint32", str)
        } else {
            this.path = path
        }
    }
    toString(){
        return this.path
    }
}

GoUint64 = class{
    constructor(str,path){
        if (typeof(str) != 'number'|| path != undefined){
            throw "GoString constructor expected number"
        }
        if (path == undefined){
            this.path=  call_go_fn("create", "uint64", str)
        } else {
            this.path = path
        }
    }
    toString(){
        return this.path
    }
}

GoUint128 = class{
    constructor(str,path){
        if (typeof(str) != 'bigint'|| path != undefined){
            throw "GoString constructor expected bigint"
        }
        if (path == undefined){
            this.path=  call_go_fn("create", "uint128", str.toString())
        } else {
            this.path = path
        }
    }
    toString(){
        return this.path
    }
}

GoFloat = class{
    constructor(str,path){
        if (typeof(str) != 'number' || path != undefined){
            throw "GoString constructor expected number"
        }
        if (path == undefined){
            this.path=  call_go_fn("create", "float", str.toString())
        } else {
            this.path = path
        }
    }
    toString(){
        return this.path
    }
}

GoFloat32 = class{
    constructor(str,path){
        if (typeof(str) != 'number' || path != undefined){
            throw "GoString constructor expected number"
        }
        if (path == undefined){
            this.path=  call_go_fn("create", "float32", str.toString())
        } else {
            this.path = path
        }
    }
    toString(){
        return this.path
    }
}

GoFloat64 = class{
    constructor(str,path){
        if (typeof(str) != 'number' || path != undefined){
            throw "GoString constructor expected number"
        }
        if (path == undefined){
            this.path=  call_go_fn("create", "float64", str.toString())
        } else {
            this.path = path
        }
    }
    toString(){
        return this.path
    }
}

GoFloat128 = class{
    constructor(str,path){
        if (typeof(str) != 'string' || path != undefined){
            throw "GoString constructor expected string"
        }
        if (path == undefined){
            this.path=  call_go_fn("create", "float128", str.toString())
        } else {
            this.path = path
        }
    }
    toString(){
        return this.path
    }
}
DB.settings = {
    is_ML_enabled: function(){
        return call_go_fn("settings","ML_enabled");
    },

    enable_ML: function(){
        return call_go_fn("settings","enable_ML");
    },

    disable_ML: function(){
        return call_go_fn("settings", "disable_ML");
    }
}

DB.adm = {
    user:user = class{
        constructor(name){
            this.name = name
            if (!call_go_fn("adm","user", "userExist",name)){
                throw "user "+name+" does not exist";
            }
        }
        addgroup(){
            return call_go_fn("adm","user","addgroup",this.name)
        }
        groupids(){
            return call_go_fn("adm","user","groupids",this.name);
        }
    },

    useradd: function(name, password){
        call_go_fn("adm","user","useradd",name,password);
        return new DB.adm.user(name)
    },

    getUser:function(name){
        return new DB.adm.user(name)
    }
}
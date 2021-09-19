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
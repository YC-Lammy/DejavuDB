fmt = function(){
    var exports = {}


    exports.Print = function(...any){
        returning_print_buffer +=  String(call_go_fn("fmt", "Print", ...any))
    }
    exports.Printf = function(...any){
        returning_print_buffer +=  String(call_go_fn("fmt", "Printf", ...any))
    }
    exports.Println = function(...any){
        returning_print_buffer +=  String(call_go_fn("fmt", "Println", ...any))
    }
    exports.Sprint = function(...any){
        return  String(call_go_fn("fmt", "Print", ...any))
    }
    exports.Sprintf = function(...any){
        return String(call_go_fn("fmt", "Printf", ...any))
    }
    exports.Sprintln = function(...any){
        return  String(call_go_fn("fmt", "Println", ...any))
    }
    return exports
}()
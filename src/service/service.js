const net = require("net")
const workerpool = require('workerpool')

const port = 5630;


const server = new net.Server()

const pool = workerpool.pool()

function handler(data){
    try{
        
        msg = JSON.parse(data.toString())

        switch (msg.name) {
            case "add_service":
                
                break;
            case "workflow":
            default:
                break;
        }

    } catch(e){
        return e
    }
}


server.listen(port)
server.once('connection',function(socket){
    socket.on('data',function(data){
        pool.exec(handler,[data]).then(function (result) { // spawn a new worker
          })
          .catch(function (err) {
            console.error(err);
          });
    })
})



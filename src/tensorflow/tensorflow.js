const net = require("net");
const tf = require("@tensorflow/tfjs-node");
const workerpool = require('workerpool')

const port = 5630;
const server = new net.Server();

const pool = workerpool.pool()

function handler(data){

}
server.listen(port)
server.on('error',function(err){
    console.error(err);
})
server.once('connection',function(socket){
    socket.on('data',function(data){
        pool.exec(handler,[data]).then(function (result) { // spawn a new worker
          })
          .catch(function (err) {
            console.error(err);
          });
    })
})
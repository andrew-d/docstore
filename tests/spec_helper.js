var path = require('path'),
    request = require('superagent'),
    spawn = require('child_process').spawn;

var config = require('./config');

var proc, exited = false;

before(function(done) {
  var exe = path.normalize(path.join(__dirname, '..', 'docstore'));

  console.log("Booting docstore server: " + exe);

  // TODO: pull these values from configuration
  proc = spawn(exe, [
    '--dbconn=:memory:',
    '--host=127.0.0.1',
    '-p', '8081',
  ], {
    stdio: ['ignore', 'ignore', 'ignore'],
  });

  proc.on('error', function(err) {
    console.log("Error spawning server process:", err);
  });

  proc.on('exit', function() {
    exited = true;
  });

  var checkServer = function() {
    console.log("Checking if server is running...");

    request.get(config.URL+'/tags')
           .end(function(err, resp) {
      if( err || !resp.ok ) {
        setTimeout(checkServer, 200);
        return;
      }

      console.log("Server started\n--------------------------------------------------\n");
      done();
    });
  };

  setTimeout(checkServer, 200);
});

after(function(done) {
  console.log("Waiting for server process to finish...");

  proc.kill('SIGINT');

  var count = 0;
  var waitExit = function() {
    if( exited ) {
      console.log("Server process finished");
      done();
      return;
    }

    if( ++count > 10 ) {
      proc.kill('SIGTERM');
    }

    setTimeout(waitExit, 100);
  };

  setTimeout(waitExit, 10);
});

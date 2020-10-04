#!/usr/bin/env node

var os           = require('os'),
fs               = require('fs'),
path             = require('path'),
exec             = require('child_process').exec,
q                = require('q'),
pkg              = require(path.join(__dirname, 'package.json'));

var program = require('nomnom').script('icns')
  .option('version', {
    abbr: 'v',
    flag: true,
    help: 'print version and exit',
    callback: function() {
      return "node-icns version " + pkg.version;
    }
  })
  .option('in', {abbr: 'i', help: 'Set the input folder name, by default it assumes the current working directory'})
  .option('out', {abbr: 'o', help: 'Set the output .icns file name, by default it will be the name of the current dir.icns'})
  .parse();

var cwd =  process.cwd();
var inName = program.in || null;
var outName = program.out || cwd.split(path.sep).pop() + '.icns';

if(inName === null) {
  var files = fs.readdirSync(cwd);
  for(var i = 0; i < files.length; i++) {
    if(fs.statSync(path.join(cwd, files[i])).isDirectory()) {
      continue;
    }
    var nameParts = files[i].split('.');
    if(nameParts.pop() === 'png') {
      inName = path.join(cwd, files[i]);
    }
  }
}

console.log('Input file:  ', inName);
console.log('Output file: ', outName);

var tmp = path.join(cwd, 'tmpIcon.iconset');
if(fs.existsSync(tmp)) {
  exec('rm -rf ' + tmp, buildIconSet);
} else {
  buildIconSet();
}

function safeExec(cmd) {
  var defer = q.defer();

  exec(cmd, function(err, stdout, stderr) {
    if(err) {
      return defer.reject(err);
    }

    return defer.resolve(stdout);
  });

  return defer.promise;
};

function buildIconSet() {
  fs.mkdirSync(tmp);
  exec('which sips', function(err, stdout, stderr) {
    if(err) {
      return console.log('ERROR:       `sips` command not on system');
    }

    var promises = [];
    var sizes = [
      16,
      32,
      128,
      256,
      512
    ]
    for(var i = 0; i < sizes.length; i++) {
      var size = sizes[i];
      promises.push((function(size) {
        return safeExec('sips -z ' + size + ' ' + size + ' ' + inName + ' --out ' + path.join(tmp, 'icon_' + size + 'x' + size + '.png'));
      })(size));
      promises.push((function(size) {
        var twiceSize = size*2;
        return safeExec('sips -z ' + twiceSize + ' ' + twiceSize + ' ' + inName + ' --out ' + path.join(tmp, 'icon_' + size + 'x' + size + '@2x.png'));
      })(size));
    }

    q.all(promises).spread(function() {
      cleanUp();
    }).fail(function(err) {
      console.log(err);
    });

  });
}

function cleanUp() {
  console.log('Running:      iconutil -c icns ' + tmp + ' -o ' + outName);
  exec('iconutil -c icns ' + tmp + ' -o ' + outName, function(err, stdout, stderr) {
    if(err) {
      return console.log(err);
    }

    exec('rm -rf ' + tmp, function(err, stdout, stderr) {
      if(err) {
        return console.log(err);
      }

      console.log('SUCCESS:      ' + outName + ' created!');
      console.log('');
    });

  });
}

#!/usr/bin/env node

"use strict"

const request = require('request'),
    path = require('path'),
    tar = require('tar'),
    zlib = require('zlib'),
    mkdirp = require('mkdirp'),
    fs = require('fs'),
    exec = require('child_process').exec;

// Mapping from Node's `process.arch` to Golang's `$GOARCH`
const ARCH_MAPPING = {
    "ia32": "32-bit",
    "x64": "64-bit",
    "arm": "arm"
};

// Mapping between Node's `process.platform` to Golang's 
const PLATFORM_MAPPING = {
    "darwin": "darwin",
    "linux": "linux",
    "win32": "windows",
    "freebsd": "freebsd"
};

function getInstallationPath(callback) {
    // `npm bin` will output the path where binary files should be installed
    exec("npm bin", function(err, stdout, stderr) {

        let dir =  null;
        if (err || stderr || !stdout || stdout.length === 0)  {

            // We couldn't infer path from `npm bin`. Let's try to get it from
            // Environment variables set by NPM when it runs.
            // npm_config_prefix points to NPM's installation directory where `bin` folder is available
            // Ex: /Users/foo/.nvm/versions/node/v4.3.0
            let env = process.env;
            if (env && env.npm_config_prefix) {
                dir = path.join(env.npm_config_prefix, "bin");
            }
        } else {
            dir = stdout.trim();
        }

        mkdirp.sync(dir);

        callback(null, dir);
    });

}

function validateConfiguration(packageJson) {
    if (!packageJson.version) {
        return "'version' property must be specified";
    }

    if (!packageJson.goBinary || typeof(packageJson.goBinary) !== "object") {
        return "'goBinary' property must be defined and be an object";
    }

    if (!packageJson.goBinary.name) {
        return "'name' property is necessary";
    }

    if (!packageJson.goBinary.path) {
        return "'path' property is necessary";
    }

    if (!packageJson.goBinary.url) {
        return "'url' property is required";
    }
}

function parsePackageJson() {
    if (!(process.arch in ARCH_MAPPING)) {
        console.error("Installation is not supported for this architecture: " + process.arch);
        return;
    }

    if (!(process.platform in PLATFORM_MAPPING)) {
        console.error("Installation is not supported for this platform: " + process.platform);
        return
    }

    const packageJsonPath = path.join(".", "package.json");
    if (!fs.existsSync(packageJsonPath)) {
        console.error("Unable to find package.json. " +
            "Please run this script at root of the package you want to be installed");
        return
    }

    let packageJson = JSON.parse(fs.readFileSync(packageJsonPath));
    let error = validateConfiguration(packageJson);
    if (error && error.length > 0) {
        console.error("Invalid package.json: " + error);
        return
    }

    // We have validated the config. It exists in all its glory
    let binName = packageJson.goBinary.name;
    let binPath = packageJson.goBinary.path;
    let url = packageJson.goBinary.url;
    let version = packageJson.version;
    if (version[0] === 'v') version = version.substr(1);  // strip the 'v' if necessary v0.0.1 => 0.0.1

    // Binary name on Windows has .exe suffix
    if (process.platform === "win32") {
        binName += ".exe"
    }

    // Interpolate variables in URL, if necessary
    url = url.replace(/{{arch}}/g, ARCH_MAPPING[process.arch]);
    url = url.replace(/{{platform}}/g, PLATFORM_MAPPING[process.platform]);
    url = url.replace(/{{version}}/g, version);
    url = url.replace(/{{bin_name}}/g, binName);

    return {
        binName: binName,
        binPath: binPath,
        url: url,
        version: version
    }
}

/**
 * Reads the configuration from application's package.json,
 * validates properties, downloads the binary, untars, and stores at
 * ./bin in the package's root. NPM already has support to install binary files
 * specific locations when invoked with "npm install -g"
 *
 *  See: https://docs.npmjs.com/files/package.json#bin
 */
const INVALID_INPUT = "Invalid inputs";
function install(callback) {
    let opts = parsePackageJson();
    if (!opts) return callback(INVALID_INPUT);

    mkdirp.sync(opts.binPath);

    console.log("Downloading and unziping build files")

    let ungz = zlib.createGunzip();
    let untar = tar.Extract({path: opts.binPath});

    untar.on('end', function() {console.log("Wio binary successfuly downloaded and extracted")});

    let req = request({uri: opts.url});
    req.on('error', callback.bind(null, "Error downloading from URL: " + opts.url));
    req.on('response', function(res) {
        if (res.statusCode !== 200) return callback("Error downloading Wio binary. HTTP Status Code: " + res.statusCode);

        req.pipe(ungz).pipe(untar);
    });
}

var myCallback = function(data) {
    console.log(data);
    process.exit(1)
};

// call the install
install(myCallback)

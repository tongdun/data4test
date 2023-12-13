const fs = require('fs');
const path = require('path');
const webpack = require('webpack');

function deleteDir(dir) {
    var files = [];
    if (fs.existsSync(dir)) {
        files = fs.readdirSync(dir);
        files.forEach(function (file, index) {
            var curDir = dir + "/" + file;
            if (fs.statSync(curDir).isDirectory()) { // recurse  
                deleteJS(curDir);
            } else { // delete file  
                fs.unlinkSync(curDir);
            }
        });
        fs.rmdirSync(dir);
    }
};

var isDev = false;
if (process.env && process.env.NODE_ENV === 'dev') {
    isDev = true;
}

//总是清除之前生成的文件
var distDir = path.join(__dirname, 'web', 'static', 'js');
deleteDir(distDir);

var config = {
    entry: {
        vendor: ['es6-promise', 'vue', 'vue-router', 'vue-i18n', 'echarts', 'lodash', 'axios', 'moment', 'iview-style', 'iview', 'highlight.js-style', 'highlight.js'],
        app: ['./web/ts/main.ts']
    },
    output: {
        filename: '[name].min.js',
        path: distDir
    },
    resolve: {
        extensions: ['.ts', '.js', '.vue'],
        alias: {
            'iview-style': 'iview/dist/styles/iview.css',
            'highlight.js-style': 'highlight.js/styles/default.css'
        }
    },
    module: {
        rules: [
            {
                test: /\.css$/,
                loader: 'style-loader!css-loader'
            },
            {
                test: /\.vue$/,
                loader: 'vue-loader',
                options: {
                    esModule: true
                }
            },
            {
                test: /\.ts$/,
                loader: 'ts-loader',
                options: {
                    appendTsSuffixTo: [/\.vue$/]
                }
            },
            {
                test: /\.(png|jpg|gif|woff|woff2|svg|eot|ttf)$/,
                loader: 'url-loader',
                options: {
                    name: './static/img/[hash].[ext]'
                }
            }
        ]
    },
    plugins: [
        new webpack.optimize.CommonsChunkPlugin({
            name: 'vendor',
            minChunks: Infinity
        })
    ]
}

if (isDev) {
    config.watch = true;
    config.devtool = '#cheap-module-eval-source-map';
    config.plugins.unshift(
        new webpack.DefinePlugin({
            'process.env': {
                NODE_ENV: '"development"'
            }
        })
    );
} else {
    config.plugins.unshift(
        new webpack.DefinePlugin({
            'process.env': {
                NODE_ENV: '"production"'
            }
        })
    );
    config.plugins.unshift(
        new webpack.optimize.UglifyJsPlugin({
            minimize: true
        })
    );

    //生产模式，删除 go/statik 目录重新生成
    deleteDir(path.join(__dirname, 'statik'));
}

module.exports = config;
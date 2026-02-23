const path = require('path');

// module.exports = {
const config = {
    // mode: 'development',
    // devtool: "inline-source-map",
    // devtool: "source-map",
    entry: {
        // https://www.youtube.com/watch?v=JcKRovPhGo8&ab_channel=Tocode
        // 14:50
        main: './static/App.js'
    },
    output: {
        // https://www.youtube.com/watch?v=JcKRovPhGo8&ab_channel=Tocode
        // 34:40
        path: path.resolve(__dirname, './static/dist'),
        // https://www.youtube.com/watch?v=JcKRovPhGo8&ab_channel=Tocode
        // 16:40 - publicPath is for dev-server
        filename: '[name].js',
        // https://www.youtube.com/watch?v=JcKRovPhGo8&ab_channel=Tocode
        // 21:50 --- dev-server keeps '/dist' im mem
        publicPath: '/dist',
        // https://docs.google.com/document/d/19rVPWxqQLlxvz9TG8mqijNPMJJb5ggwhMpMzVhRZMU4/edit
        clean: true,
    },
    devServer: {
        port: 3001,
        open: true,
        static: {directory: path.join(__dirname)}
    },
    module: {
        rules: [
            {
                test: /\.js$/,
                exclude: /node_modules/,
                use: {
                    loader: "babel-loader",
                    options: {
                        presets: ['@babel/preset-env']
                    }
                }
            },
            {
                test: /\.css$/,
                use: [
                    "style-loader",
                    "css-loader"
                ]
            }
        ]
    }
};

// https://webpack.js.org/configuration/mode/
// If you want to change the behavior according to the mode variable inside the webpack.config.js,
// you have to export a function instead of an object:

module.exports = (env, argv) => {

    // === panedrone: 'argv.mode' is from '> webpack mode development'
    if (argv.mode === 'development') {
        // config.devtool = 'source-map';
        config.devtool = 'inline-source-map'
    }

    if (argv.mode === 'production') {
        //...
    }

    return config;
};
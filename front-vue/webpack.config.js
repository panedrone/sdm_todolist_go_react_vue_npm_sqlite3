const path = require('path');
//const VueLoaderPlugin = require('vue-loader/lib/plugin')
// https://github.com/vuejs/vue-loader/issues/1688
const {VueLoaderPlugin} = require('vue-loader')

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
    resolve: { // package.json.PANEDRONE.vue_2.7
        alias: {
            'vue$': 'vue/dist/vue.esm.js'
        }
    },
    output: {
        // https://www.youtube.com/watch?v=JcKRovPhGo8&ab_channel=Tocode
        // 34:40
        path: path.resolve(__dirname, './static/dist'),
        // https://www.youtube.com/watch?v=JcKRovPhGo8&ab_channel=Tocode
        // 16:40
        filename: '[name].js',
        // publicPath is required for dev-server
        // https://www.youtube.com/watch?v=JcKRovPhGo8&ab_channel=Tocode
        // 21:50 --- dev server stores "dist" folder in memory
        // '/dist' without leading "." (===panedrone: http://localhost:3001/dist/main.js)
        publicPath: '/dist',
        // https://docs.google.com/document/d/19rVPWxqQLlxvz9TG8mqijNPMJJb5ggwhMpMzVhRZMU4/edit
        clean: true
    },
    devServer: {
        port: 3001,
        open: true,

        // https://www.youtube.com/watch?v=JcKRovPhGo8&ab_channel=Tocode
        //  коментах
        /*
            @nikitaberdyev
            2 years ago
            столкнулся с несколькими проблемами и вот как их решил:
            просто написать "overlay:true" больше нельзя, пишем: "client: {overlay: true}"
            также в devServer нужно указать директорию, так как по умолчанию он берет файлы из ./public
                поэтому в devServer прописываем "static: {directory: path.join(__dirname)}"
         */

        // === panedrone: а оно и правильно, ему нужно скопировать 'index.html' к себе на http://localhost:3001

        static: {directory: path.join(__dirname)}
    },
    module: {
        rules: [
            {
                test: /\.vue$/,
                use: [
                    "vue-loader"
                ]
            },
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
            // {
            //     test: /\.css$/,
            //     use: [
            //         "style-loader",
            //         "css-loader"
            //     ]
            // }
        ]
    },
    plugins: [
        new VueLoaderPlugin()
    ]
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
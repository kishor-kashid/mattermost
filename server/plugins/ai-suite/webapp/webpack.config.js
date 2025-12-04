const path = require('path');
const webpack = require('webpack');
const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const manifest = require('../plugin.json');

const DIST_PATH = path.resolve(__dirname, 'dist');

module.exports = (env = {}, argv = {}) => {
    const mode = argv.mode || 'production';

    return {
        mode,
        devtool: mode === 'production' ? 'source-map' : 'inline-source-map',
        entry: './src/index.tsx',
        output: {
            path: DIST_PATH,
            filename: 'main.js',
            publicPath: '/',
        },
        resolve: {
            extensions: ['.ts', '.tsx', '.js', '.jsx'],
        },
        module: {
            rules: [
                {
                    test: /\.tsx?$/,
                    use: 'ts-loader',
                    exclude: /node_modules/,
                },
                {
                    test: /\.css$/,
                    use: [MiniCssExtractPlugin.loader, 'css-loader'],
                },
            ],
        },
        externals: {
            react: 'React',
            'react-dom': 'ReactDOM',
        },
        plugins: [
            new MiniCssExtractPlugin({
                filename: 'main.css',
            }),
            new ForkTsCheckerWebpackPlugin(),
            new webpack.DefinePlugin({
                'process.env.PLUGIN_ID': JSON.stringify(manifest.id),
                'process.env.PLUGIN_VERSION': JSON.stringify(manifest.version),
            }),
        ],
        optimization: {
            minimize: mode === 'production',
        },
    };
};


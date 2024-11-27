const { defineConfig } = require("@vue/cli-service");
const path = require("path");
const CompressionWebpackPlugin = require("compression-webpack-plugin");
const UglifyJsPlugin = require('uglifyjs-webpack-plugin');


function resolve(dir) {
	return path.join(__dirname, dir);
}

module.exports = defineConfig({
	publicPath:'/',
	transpileDependencies: true,
	productionSourceMap: false,
	outputDir: process.env.outputDir,
	configureWebpack: (config) => {
		const plugins = [
			// build 删除日志
			// new UglifyJsPlugin({
			// 	uglifyOptions:{
			// 		compress:{
			// 			drop_console:true
			// 		}
			// 	}
			// })
		];
		// start 生成 gzip 压缩文件
		// config.optimization.splitChunks({chunks: 'all'});
		// config.plugins.delete('prefetch');
		plugins.push(
			new CompressionWebpackPlugin({
				// filename: '[path].gz[query]', // 目标资源名称
				// algorithm: 'gzip',
				// test: productionGzipExtensions, // 处理所有匹配此 {RegExp} 的资源
				// threshold: 10240, // 只处理比这个值大的资源。按字节计算(楼主设置10K以上进行压缩)
				// minRatio: 0.8 // 只有压缩率比这个值小的资源才会被处理
				algorithm: "gzip",
				// 匹配压缩文件
				test: /\.js$|\.css$/,
				// 对于大于10k压缩
				threshold: 10240,
			})
		);
		// End 生成 gzip 压缩文件
		config.plugins = [...config.plugins, ...plugins];
	},
	chainWebpack: (config) => {
		config.plugin("html").tap((args) => {
			args[0].title = "AskOnce";
			return args;
		});
		//最小化代码
		config.optimization.minimize(true);
		//分割代码
		config.optimization.splitChunks({ chunks: "all" });
		//默认开启prefetch(预先加载模块)，提前获取用户未来可能会访问的内容 在首屏会把这十几个路由文件，都一口气下载了 所以我们要关闭这个功能模块
		config.plugins.delete("prefetch");
		config.module.rules.delete("svg"); // 重点:删除默认配置中处理svg
		config.module
			.rule("svg-sprite-loader") // rule 匹配规则
			.test(/\.svg$/) // 用正则匹配 文件
			.include // 包含
			.add(resolve("src/icons")) // 处理svg目录
			.end()
			// eslint-disable-next-line no-irregular-whitespace
			.use("svg-sprite-loader") //配置loader  use() 使用哪个loader
			.loader("svg-sprite-loader") // 加载loader
			.options({
				// [name] 变量。一般表示匹配到的文件名 xxx.svg
				// eslint-disable-next-line no-irregular-whitespace
				// 注意： symbolId  在  <use xlink:href="#dl-icon-svg文件名" />
				symbolId: "icon-[name]", // 将所有的.svg 集成到 symbol中，当使用 类名 icon-文件名
			});
		config.module
			.rule("md")
			.test(/\.md$/)
			.use("text-loader")
			.loader("text-loader")
			.end();
	},
	devServer: {
		proxy: {
			"/serverApi": {
				target: "https://gateway.atomecho.cn",
				pathRewrite: {
					"^/serverApi": "/",
				},
				changeOrigin: true,
				onProxyRes(proxyRes) {
					// 设置响应头，允许接收流数据
					proxyRes.headers["transfer-encoding"] = "chunked";
					proxyRes.headers["content-type"] = "application/octet-stream";
				},
			},
		},
	},
});

module.exports = {
	root: true,
	env: {
		node: true,
	},
	extends: ["plugin:vue/vue3-essential", "eslint:recommended"],
	parserOptions: {
		parser: "@babel/eslint-parser",
	},
	rules: {
		"vue/multi-word-component-names": "off",
		"no-console": "off",
		"no-debugger": "off",
		"no-unused-vars": "off",
	},
	globals: {
		Aliplayer: true,
	},
};

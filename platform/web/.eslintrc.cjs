module.exports = {
	root: true,
	env: { browser: true, es2020: true },
	extends: [
		"eslint:recommended",
		"plugin:@typescript-eslint/recommended",
		"plugin:react-hooks/recommended",
		"prettier",
		"plugin:prettier/recommended"
	],
	ignorePatterns: ["dist", ".eslintrc.cjs"],
	parser: "@typescript-eslint/parser",
	plugins: ["react-refresh", "prettier"],
	rules: {
		"react-refresh/only-export-components": [
			"warn",
			{ allowConstantExport: true }
		],
		"prettier/prettier": "error",
		"react/react-in-jsx-scope": "off"
	}
};

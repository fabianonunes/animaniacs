const { readFileSync } = require('fs')
const { resolve } = require('path')

exports.yakko = readFileSync(resolve(__dirname, './art/YAKKO'))
exports.wakko = readFileSync(resolve(__dirname, './art/WAKKO'))
exports.dot = readFileSync(resolve(__dirname, './art/DOT'))

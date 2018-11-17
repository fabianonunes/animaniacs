const { sample } = require('lodash')
const { yakko, wakko, dot } = require('./animaniacs')

const chars = { yakko, wakko, dot }

const randomChar = (req, res) => { 
  res.setHeader('Content-Type', 'text/plain')
  return sample(chars)
}
const char = (req, res) => { 
  res.setHeader('Content-Type', 'text/plain')
  return chars[process.env.CHAR]
}

exports.default = process.env.CHAR ? char : randomChar

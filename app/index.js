const { sample } = require('lodash')
const query = require('micro-query')
const sleep = require('then-sleep')
const { yakko, wakko, dot } = require('./animaniacs')

const chars = { yakko, wakko, dot }

const randomChar = (req, res) => {
  res.setHeader('Content-Type', 'text/plain; charset=utf-8')
  return sample(chars)
}

const char = async (req, res) => {
  res.setHeader('Content-Type', 'text/plain; charset=utf-8')

  const params = query(req)
  if (params.sleep) {
    await sleep(parseInt(params.sleep, 10))
  }

  return chars[process.env.CHAR]
}

exports.default = process.env.CHAR ? char : randomChar

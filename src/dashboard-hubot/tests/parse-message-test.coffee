chai = require('chai')
expect = chai.expect
parser = require '../scripts/parse-message.coffee'

describe 'parser', ->
  describe '.parseName(string)', ->
    it 'should return the first string of two strings separated by one space as name', ->
      expect(parser.parseName('hello world')).to.equal('hello')

  describe '.parseUrl(string)', ->
    it 'should return the second string of two strings separated by one space as url', ->
      expect(parser.parseUrl('hello world')).to.equal('world')

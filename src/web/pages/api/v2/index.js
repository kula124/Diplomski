const db = require('../../../db')

export default async function handler(req, res) {
  if (req.method === "POST") {
    const { key, hash, paid } = req.body
    if (!(!!key && !!hash)) {
      return res.status(400).setHeader('Content-Type', 'application/json').json({ msg: "Wrong request format" })
    }
    const { error } = await db.postKey({ key, hash, paid })
    if (error) {
      if (error.code === "23505") { // duplicate key error
        return res.status(200).setHeader('Content-Type', 'application/json').json({ msg: "request accepted" })
      }
      return res.status(400).setHeader('Content-Type', 'application/json').json({ msg: "Key request declined" })
    }
    return res.status(200).setHeader('Content-Type', 'application/json').json({ msg: "request accepted" })
  }
  res.status(404).setHeader("Content-Type", "application/json").json({ message: "Endpoint or METHOD incorrect" })
}
const db = require('../../../db')

export default async function handler(req, res) {
  if (req.method === "GET") {
    console.log(req.headers)
    const authHeader = req.headers['authorization']
    if (!authHeader) {
      return res.status(401).json({ message: "Unauthorized access! This attempt will be logged" })
    }
    const token = authHeader.split(' ')[1]
    const userInfo = await db.userInfo(token)
      .catch(err => {
        console.log("ERROR M8")
        return res.status(400).json({ message: "failed to get user info" })
      })
    return res.status(200).json(userInfo || {})
  }
  res.status(404).json({ message: "method or endpoint incorrect" })
}
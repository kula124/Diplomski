const db = require('../../../db')

export default async function handler(req, res) {
  if (req.method === "POST") {
    const { email, password } = req.body
    const { user, session } = await db.auth(email, password)
    // console.log(user, session)
    return res.status(200).setHeader('Content-Type', 'application/json').json({ user, session })
  }
  res.status(404).setHeader("Content-Type", "application/json").json({ message: "Endpoint or METHOD incorrect" })
}
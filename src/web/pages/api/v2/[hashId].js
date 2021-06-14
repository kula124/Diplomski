const db = require('../../../db')

export default async function handler(req, res) {
  if (req.method === "GET") {
    const { data, error } = await db.getPaidKeyByHash(req.query.hashId)
    if (error) {
      return res.status(200).setHeader("Content-Type", "application/json").json({ ...error })
    }
    return res.status(200).setHeader("Content-Type", "application/json").json({ ...data })
  }
  res.status(404).setHeader("Content-Type", "application/json").json({ message: "Endpoint or METHOD incorrect" })
}

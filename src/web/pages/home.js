import Head from 'next/head'

import HomePage from '../components/HomePage'
import Protector from '../components/ProtectorComponent'

export default function Home() {
  return (
    <div>
      <Head>
        <title>RW CnC</title>
        <meta name="description" content="cool page brah" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Protector shouldBeLoggedIn={true} goTo={"/"}>
        <HomePage />
      </Protector>
    </div>
  )
}

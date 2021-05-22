import Head from 'next/head'

import Protector from '../components/ProtectorComponent'
import V1 from '../components/v1Page'

export default function Home() {
  return (
    <div>
      <Head>
        <title>RW CnC</title>
        <meta name="description" content="cool page brah" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Protector shouldBeLoggedIn={true} goTo={"/"}>
        <V1 />
      </Protector>
    </div>
  )
}

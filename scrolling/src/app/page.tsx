'use client'
import Image from "next/image";
import { useRef } from "react";
import gsap from "gsap";
import { useGSAP } from "@gsap/react";

export default function Home() {
  const container = useRef(null)
  const { contextSafe } = useGSAP()
  gsap.registerPlugin(useGSAP)

  useGSAP(() => {
    gsap.to('.box', { rotation: '+=360', duration: 4 })
  }, {
    scope: container
  })

  const onEnter = contextSafe(({ currentTarget }: React.MouseEvent) => {
    gsap.to(currentTarget, { rotation: '+=360' })
  })

  return (
    <main ref={container} className="flex items-center justify-center min-h-screen m-0 overflow-hidden">
      <div className="box bg-green-600">Rotation</div>
      <div className="boxClick bg-blue-700" onClick={onEnter}>
        Click Me
      </div>
    </main>
  );
}

'use client'
import SVGComponent from "../ui/SVGComponent";
import { useRef } from "react";
import gsap from "gsap";
import { useGSAP } from "@gsap/react";
// Remove the import statement for MotionPathPlugin
import { ScrollTrigger, SVGPathValue } from 'gsap/ScrollTrigger';
import MotionPathPlugin from "gsap/MotionPathPlugin";

gsap.registerPlugin(useGSAP, ScrollTrigger, MotionPathPlugin);

export default function Page() {
  const container = useRef(null)


  useGSAP(
    () => {
      gsap.to(".motionSVG", {
        scrollTrigger: {
          trigger: ".motionPath",
          start: "top center",
          end: () => "+=" + (document.querySelector(".motionPath")?.getBoundingClientRect()?.height || 0),
          scrub: 0.5,
          markers: true,
        },
        ease: pathEase(".motionPath"), // a custom ease that helps keep the tractor centered
        immediateRender: true,
        motionPath: {
          path: ".motionPath",
          align: ".motionPath",
          alignOrigin: [0.5, 0.5],
          autoRotate: 90,
        }
      });
    },
    { scope: container }
  );

  function pathEase(path: string, config: { axis?: string, precision?: number, smooth?: boolean | number } = {}) {
    let axis = config.axis || "y",
      precision = config.precision || 1,
      rawPath = MotionPathPlugin?.cacheRawPathMeasurements(MotionPathPlugin.getRawPath(gsap.utils.toArray(path)[0] as SVGPathValue), Math.round(precision * 12)),
      useX = axis === "x",
      start = rawPath?.[0][useX ? 0 : 1],
      end = rawPath?.[rawPath?.length - 1][rawPath[rawPath?.length - 1].length - (useX ? 2 : 1)],
      range = end - start,
      l = Math.round(precision * 200),
      inc = 1 / l,
      positions = [0],
      a = [],
      minIndex = 0,
      smooth = [0],
      minChange = (1 / l) * 0.6,
      smoothRange = config.smooth === true ? 7 : Math.round(config.smooth || 0) || 0,
      fullSmoothRange = smoothRange * 2,
      getClosest = (p: number) => {
        while (positions[minIndex] <= p && minIndex++ < l) { }
        a.push((p - positions[minIndex - 1]) / (positions[minIndex] - positions[minIndex - 1]) * inc + minIndex * inc);
        smoothRange && a.length > smoothRange && (a[a.length - 1] - a[a.length - 2] < minChange) && smooth.push(a.length - smoothRange);
      },
      i = 1;
    for (; i < l; i++) {
      positions[i] = (MotionPathPlugin?.getPositionOnPath(rawPath, i / l)[axis] - start) / range;
    }
    positions[l] = 1;
    for (i = 0; i < l; i++) {
      getClosest(i / l);
    }
    a.push(1); // must end at 1.
    if (smoothRange) { // smooth at the necessary indexes where a small difference was sensed. Make it a linear change over the course of the fullSmoothRange
      smooth.push(l - fullSmoothRange + 1);
      smooth.forEach(i => {
        let start = a[i],
          j = Math.min(i + fullSmoothRange, l),
          inc = (a[j] - start) / (j - i),
          c = 1;
        i++;
        for (; i < j; i++) {
          a[i] = start + inc * c++;
        }
      });
    }
    return (p: number) => {
      let i = p * l,
        s = a[i | 0];
      return i ? s + (a[Math.ceil(i)] - s) * (i % 1) : 0;
    }
  }



  return (
    <div ref={container}>
      <SVGComponent />
    </div >
  )
}
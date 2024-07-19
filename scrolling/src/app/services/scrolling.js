// for another helper function that calculates the exact progress value along a motion path where it'll hit the center of the provided target on the given axis ("y" by default), see https://codepen.io/GreenSock/pen/BaPdrKM
// see also: https://codepen.io/GreenSock/pen/mdgmawg
gsap.registerPlugin(MotionPathPlugin, ScrollTrigger);

gsap.set("#motionSVG", { scale: 0.7, autoAlpha: 1 });
gsap.set("#tractor", {transformOrigin: "50% 50%"});
let rotateTo = gsap.quickTo("#tractor", "rotation"),
    prevDirection = 0;

gsap.to("#motionSVG", {
  scrollTrigger: {
    trigger: "#motionPath",
    start: "top center",
    end: () => "+=" + document.querySelector("#motionPath").getBoundingClientRect().height,
    scrub: 0.5,
    markers: true,
    onUpdate: self => {
      if (prevDirection !== self.direction) { // only run this when we're changing direction
        rotateTo(self.direction === 1 ? 0 : -180);
        prevDirection = self.direction;
      }
    }
  },
  ease: pathEase("#motionPath"), // a custom ease that helps keep the tractor centered
  immediateRender: true,
  motionPath: {
    path: "#motionPath",
    align: "#motionPath",
    alignOrigin: [0.5, 0.5],
    autoRotate: 90,
  }
});




/* 
Helper function that returns an ease that bends time to ensure the target moves on the y axis in a relatively steady fashion in relation to the viewport (assuming the progress of the tween is linked linearly to the scroll position). Requires MotionPathPlugin of course.
You can optionally pass in a config option with any of these properties: 
  - smooth: if true, the target can drift slightly in order to smooth out the movement. This is especially useful if the path curves backwards at times. It prevents super-fast motions at that point. You can define it as a number (defaults to 7) indicating how much to smooth it.
  - precision: number (defaults to 1) controlling the sampling size along the path. The higher the precision, the more accurate but the more processing.
  - axis: "y" or "x" ("y" by default)
*/ 
function pathEase(path, config={}) {
  let axis = config.axis || "y",
      precision = config.precision || 1,
      rawPath = MotionPathPlugin.cacheRawPathMeasurements(MotionPathPlugin.getRawPath(gsap.utils.toArray(path)[0]), Math.round(precision * 12)),
			useX = axis === "x",
			start = rawPath[0][useX ? 0 : 1],
			end = rawPath[rawPath.length - 1][rawPath[rawPath.length-1].length - (useX ? 2 : 1)],
			range = end - start,
			l = Math.round(precision * 200),
			inc = 1 / l,
			positions = [0],
			a = [],
			minIndex = 0,
      smooth = [0],
      minChange = (1 / l) * 0.6,
      smoothRange = config.smooth === true ? 7 : Math.round(config.smooth) || 0,
      fullSmoothRange = smoothRange * 2,
			getClosest = p => {
				while (positions[minIndex] <= p && minIndex++ < l) { }
				a.push((p - positions[minIndex-1]) / (positions[minIndex] - positions[minIndex - 1]) * inc + minIndex * inc);
        smoothRange && a.length > smoothRange && (a[a.length - 1] - a[a.length - 2] < minChange) && smooth.push(a.length - smoothRange);
			},
			i = 1;
  for (; i < l; i++) {
    positions[i] = (MotionPathPlugin.getPositionOnPath(rawPath, i / l)[axis] - start) / range;
  }
  positions[l] = 1;
  for (i = 0; i < l; i++) {
    getClosest(i / l);
  }
  a.push(1); // must end at 1.
  if (smoothRange) { // smooth at the necessary indexes where a small difference was sensed. Make it a linear change over the course of the fullSmoothRange
    smooth.push(l-fullSmoothRange+1);
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
  return p => {
    let i = p * l,
        s = a[i | 0];
    return i ? s + (a[Math.ceil(i)] - s) * (i % 1) : 0;
  }
}
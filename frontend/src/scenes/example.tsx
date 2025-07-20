import { Rect, SVG, makeScene2D } from '@motion-canvas/2d';
import { createRef, createSignal, createComputed, all, waitFor, easeInOutCubic } from '@motion-canvas/core';
import hospitalSvg from '/public/hospital.svg?raw';

export default makeScene2D(function* (view) {
  view.fill('#000000'); // Set the background to black

  // Create signals for dynamic properties
  const iconScale = createSignal(0);

  // Create a parent rectangle that is center-anchored
  const parentRect = createRef<Rect>();
  const hospitalIcon = createRef<SVG>();

  view.add(
    <Rect 
      ref={parentRect}
      width={() => view.width() * 0.6}
      height={() => view.height() * 0.6}
      fill={null} // No fill for the parent rect
    >
      <SVG 
        ref={hospitalIcon}
        svg={hospitalSvg}
        width={() => parentRect().width() * 0.6}
        height={() => parentRect().height() * 0.6}
        scale={() => iconScale()} // Scale based on iconScale signal
      />
    </Rect>
  );

  // Animate the scale of the icon from 0 to 1
  yield* all(
    iconScale(1, 1.2, easeInOutCubic), // Scale in the hospital icon
    waitFor(0.5) // Optional wait before the next action
  );
});

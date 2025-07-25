import {Node, makeScene2D, SVG, Path} from '@motion-canvas/2d';
import {createRef, waitFor, all} from '@motion-canvas/core';
import beach from '/public/beach_1.svg?raw';
import sun from '/public/sun_1.svg?raw';

export default makeScene2D(function* (view) {
  // Set background color
  view.fill('#000000');

  // Create references for SVGs and container
  const containerRef = createRef<Node>();
  const beachRef = createRef<SVG>();
  const sunRef = createRef<SVG>();

  view.add(
    <Node ref={containerRef}>
      <SVG ref={beachRef} svg={beach} size={440} y={() => containerRef().height() / 2 - 200} />
      <SVG ref={sunRef} svg={sun} size={140} y={() => -containerRef().height() / 2 + 120} />
    </Node>
  );

  // Animate beach and sun
  yield* all(
    beachRef().scale(1, 0.5), 
    sunRef().scale(1, 0.5),
    beachRef().opacity(1, 0.5),
    sunRef().opacity(1, 0.5)
  );

  // Animate SVG path children to fill white
  for (const child of beachRef().children()[0].children()) {
    if (child instanceof Path) {
      yield* child.fill('white', 1);
    }
  }
  for (const child of sunRef().children()[0].children()) {
    if (child instanceof Path) {
      yield* child.fill('white', 1);
    }
  }
});
import {
  Rect,
  Node,
  makeScene2D,
  SVG,
  Path,
} from '@motion-canvas/2d';
import {
  createRef,
  all,
  waitFor,
} from '@motion-canvas/core';
import moon from '/public/moon.svg?raw';

export default makeScene2D(function* (view) {
  // Set background color
  view.fill('#000000');

  // Reference for the moon SVG and container
  const moonRef = createRef<SVG>();
  const mainRectRef = createRef<Rect>();

  view.add(
    <Node>
      <Rect
        ref={mainRectRef}
        width={600}
        height={600}
        fill={null}
        x={0}
        y={0}
      >
        <SVG
          ref={moonRef}
          svg={moon}
          size={300}
          opacity={0}
          scale={0}
        />
      </Rect>
    </Node>
  );

  // Set the fill color of every Path child in the moon SVG to pale yellow
  for (const child of moonRef().children()[0].children()) {
    if (child instanceof Path) {
      yield* child.fill('#fffbe0', 0);
    }
  }

  // Animate the moon SVG scale and opacity
  yield* all(
    moonRef().scale(1, 0.5), // Animate scale to 1
    moonRef().opacity(1, 0.5) // Animate opacity to 1
  );
});

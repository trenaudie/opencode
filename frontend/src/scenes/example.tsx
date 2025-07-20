import {
  Rect,
  Node,
  SVG,
  Path,
  makeScene2D,
} from '@motion-canvas/2d';
import {
  createRef,
  createSignal,
  all,
  waitFor,
  Vector2,
} from '@motion-canvas/core';
import person from '/public/person_1.svg?raw';
import sun from '/public/sun_1.svg?raw'; // Import the sun SVG

const numRows = 5;
const baseWidth = 560;
const widthStep = 80;
const boxHeight = 48;
const gap = 8;
const svgPersonSize = 110;
const gapBetweenSandAndPerson = 56;

export default makeScene2D(function* (view) {
  // Background
  view.fill('#000000');

  // Sand rectangles setup
  const rectRefs = Array.from({ length: numRows }, () => createRef<Rect>());
  const rectSignals = Array.from({ length: numRows }, () => createSignal(0));
  const container = createRef<Node>();

  // Person SVG setup
  const personRef = createRef<SVG>();
  const personOpacity = createSignal(0);
  const personScale = createSignal(0);

  // Sun SVG setup
  const sunRef = createRef<SVG>();
  const sunOpacity = createSignal(0);
  const sunScale = createSignal(0);

  // Add sand, person, and sun to scene
  view.add(
    <>
      <Node ref={container}>
        {Array.from({ length: numRows }, (_, i) => {
          const width = baseWidth - i * widthStep;
          const y = () => (i * (boxHeight + gap)) - (numRows * (boxHeight + gap) / 2) + (boxHeight / 2);
          return (
            <Rect
              ref={rectRefs[i]}
              width={width}
              height={boxHeight}
              fill={'#F6C97F'}
              y={y}
              opacity={rectSignals[i]}
            />
          );
        })}
      </Node>
      {/* Person SVG (standing above the sand) */}
      <SVG
        ref={personRef}
        svg={person}
        size={svgPersonSize}
        opacity={personOpacity}
        scale={personScale}
        position={() => {
          // Get absolute top of highest (first) sand rectangle
          const sandAbs = rectRefs[0]().absolutePosition();
          // Convert to view's local coordinates
          const localSand = sandAbs.transformAsPoint(view.worldToLocal());
          // Place above (minus y), leave a relative gap
          return new Vector2([localSand.x, localSand.y - boxHeight / 2 - gapBetweenSandAndPerson]);
        }}
      />
      {/* Sun SVG in the top-left corner */}
      <SVG
        ref={sunRef}
        svg={sun}
        size={80} // Size of the sun SVG
        opacity={sunOpacity}
        scale={sunScale}
        position={() => new Vector2([0, 0])} // Top-left corner, offset can be adjusted
      />
    </>
  );

  // Animate rectangles (sand)
  for (let i = 0; i < numRows; i++) {
    yield* all(
      rectSignals[i](1, 0.5),
      waitFor(0.14)
    );
  }

  // Animate the person SVG scaling/fading in
  yield* all(
    personOpacity(1, 0.7),
    personScale(1, 0.7)
  );

  // Animate the fill of each path in the person SVG (to white for clarity)
  for (const child of personRef().children()[0].children()) {
    if (child instanceof Path) {
      yield* child.fill('white', 0.8);
    }
  }

  // Animate sun SVG fading in and scaling from 0 to 1
  yield* all(
    sunOpacity(1, 0.7),
    sunScale(1, 0.7)
  );
  for (const child of sunRef().children()[0].children()) {
    if (child instanceof Path) {
      yield* child.fill('yellow', 1); // Fill sun paths with yellow
    }
  }
});

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
import sun from '/public/sun_1.svg?raw'; // Import the sun SVG

const numRows = 5;
const baseWidth = 560;
const widthStep = 80;
const boxHeight = 48;
const gap = 8;

export default makeScene2D(function* (view) {
  // Background
  view.fill('#000000');

  // Sand rectangles setup
  const rectRefs = Array.from({ length: numRows }, () => createRef<Rect>());
  const rectSignals = Array.from({ length: numRows }, () => createSignal(0));
  const container = createRef<Node>();

  // Sun SVG setup
  const sunRef = createRef<SVG>();
  const sunOpacity = createSignal(0);
  const sunScale = createSignal(0);

  // Add sand and sun to scene
  view.add(
    <>
      <Node ref={container}>
        {Array.from({ length: numRows }, (_, i) => {
          // Flip: largest at top (i=0), smallest at bottom (i=numRows-1)
          const width = baseWidth - (numRows - 1 - i) * widthStep;
          // Pyramid "stands" on its tip at bottom
          const y = () => (-(i * (boxHeight + gap)) + (numRows * (boxHeight + gap) / 2) - (boxHeight / 2));
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

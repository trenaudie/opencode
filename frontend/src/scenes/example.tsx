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
import franceSVG from '/public/france_1.svg?raw';

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

  // France SVG setup
  const franceRef = createRef<SVG>();
  const franceOpacity = createSignal(0);
  const franceScale = createSignal(0.6);

  // Add sand, france, and sun to scene
  view.add(
    <>
      {/* France SVG at the top, centered horizontally */}
      <Node>
        <SVG
          ref={franceRef}
          svg={franceSVG}
          // Make width relative to view width, or a fixed max
          size={() => Math.min(view.width() * 0.22, 210)}
          opacity={franceOpacity}
          scale={franceScale}
          position={() => {
            // Top center: x=0, y = small offset downward from top
            const svgSize = Math.min(view.width() * 0.22, 210);
            return new Vector2([0, -view.height()/2 + svgSize/2 + 28]);
          }}
        />
      </Node>
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

  // Animate France SVG: fade and scale in, then fill its paths
  yield* all(
    franceOpacity(1, 0.6),
    franceScale(1, 0.6)
  );
  for (const child of franceRef().children()[0].children()) {
    if (child instanceof Path) {
      yield* child.fill('white', 0.5);
    }
  }

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

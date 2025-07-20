import {
  Rect,
  Node,
  makeScene2D,
} from '@motion-canvas/2d';
import {
  createSignal,
  createRef,
  all,
  waitFor,
  linear
} from '@motion-canvas/core';

const numRows = 5;
const baseWidth = 560;
const widthStep = 80;
const boxHeight = 48;
const gap = 8;

export default makeScene2D(function* (view) {
  // Set background color
  view.fill('#000000');

  // Refs for rectangles and signals for opacity
  const rectRefs = Array.from({ length: numRows }, () => createRef<Rect>());
  const rectSignals = Array.from({ length: numRows }, () => createSignal(0));

  // Main container Node to stack rectangles
  const container = createRef<Node>();
  view.add(
    <Node ref={container}>
      {Array.from({ length: numRows }, (_, i) => {
        const width = baseWidth - i * widthStep;
        const y = () => (i * (boxHeight + gap)) - (numRows * (boxHeight + gap) / 2) + (boxHeight / 2);

        return (
          <Rect
            ref={rectRefs[i]}
            width={width}
            height={boxHeight}
            fill={'#6cf1c2'}
            y={y}
            opacity={rectSignals[i]} // opacity controlled by signal
          />
        );
      })}
    </Node>
  );

  // Fade in rectangles one by one
  for (let i = 0; i < numRows; i++) {
    yield* all(
      rectSignals[i](1, 0.5, linear), // Fade to opaque
      waitFor(0.2) // Wait before next rectangle
    );
  }
});

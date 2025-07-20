import { Rect, makeScene2D } from '@motion-canvas/2d';
import { createSignal, all, waitFor } from '@motion-canvas/core';

export default makeScene2D(function* (view) {
  // Set the background color of the view
  view.fill('#000000');

  // Create a signal to hold the scale value of the square
  const scale = createSignal(0.4); // Start small

  // Add the square to the scene
  const squareRef = createRef<Rect>();
  view.add(
    <Rect
      ref={squareRef}
      x={() => 0} // Centered in x
      y={() => 0} // Centered in y
      width={() => 240} // Fixed width
      height={() => 240} // Fixed height
      fill={'#3485E7'} // Square fill color
      stroke={'#FFFFFF'} // Square stroke color
      lineWidth={8} // Stroke width
      scale={scale} // Reactive scale
    />
  );

  // Wait a moment to ensure everything is rendered
  yield* waitFor(0.2);

  // Animate: scale up the square from 0.4 to 1.0
  yield* all(
    squareRef().scale(0.4, 0), // Start small
    squareRef().scale(1.0, 1.2) // Scale up over 1.2 seconds
  );
});

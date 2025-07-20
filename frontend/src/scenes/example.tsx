import { Rect, makeScene2D } from '@motion-canvas/2d';
import { createRef } from '@motion-canvas/core';

export default makeScene2D(function* (view) {
  // Set the background color
  view.fill('#000000');

  // Create references for the rectangle
  const rectRef = createRef<Rect>();

  // Reactive dimensions
  const rectWidth = () => view.width() * 0.8;  // 80% of the view's width
  const rectHeight = () => view.height() * 0.2; // 20% of the view's height

  // Add the rectangle to the view
  view.add(
    <Rect
      ref={rectRef}
      width={rectWidth}
      height={rectHeight}
      fill={'#3498db'} // Filling with a visible blue color
      position={[0, 0]} // Centered in the view
    />
  );
});

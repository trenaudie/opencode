import { makeScene2D, Rect, Node, SVG, Path } from '@motion-canvas/2d';
import {  createRef, createSignal, all, waitFor , easeOutBack} from '@motion-canvas/core';
import hospitalSVG from '/public/hospital_1.svg?raw';

export default makeScene2D(function* (view) {
  // Set background color
  view.fill('#000000');

  // Create a reference for the container and the SVG
  const container = createRef<Rect>();
  const svgRef = createRef<SVG>();

  // Initialize the container
  view.add(
    <Rect
      ref={container}
      width={400}
      height={400}
    >
      <SVG
        ref={svgRef}
        svg={hospitalSVG}
        size={300} // Set base size
        scale={0} // Start scaled down
        opacity={1}
      />
    </Rect>
  );

  // Animation: Pop in the hospital icon
  yield* svgRef().scale(1, 0.9, easeOutBack); // Scale up the SVG from zero

  // Animate the paths filling from gray to white
  // Correct hierarchy: paths are usually under Node in SVG
  const nodeChildren = svgRef().children()[0]?.children() || [];
  for (const child of nodeChildren) {
    if (child instanceof Path) {
      yield* child.fill('gray', 0); // Set initial fill to gray
      yield* child.fill('white', 0.3); // Animate fill to white
      yield* waitFor(0.1); // Wait slightly between fills
    }
  }
});

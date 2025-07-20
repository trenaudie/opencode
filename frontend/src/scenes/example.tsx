import {
  Rect,
  Circle,
  Node,
  SVG,
  Path,
  makeScene2D,
} from '@motion-canvas/2d';
import {
  createRef,
  all,
  waitFor,
} from '@motion-canvas/core';
import beach from '/public/beach.svg?raw'; // Import the Beach SVG

export default makeScene2D(function* (view) {
  // Set background color
  view.fill('#000000');

  // Reference for the beach SVG
  const beachRef = createRef<SVG>();
  const sunRef = createRef<Circle>();
  const beachRectRef = createRef<Rect>();

  // Add the beach SVG centered in the scene
  view.add(
    <Node>
      <SVG
        ref={beachRef}
        svg={beach}
        size={340}
        opacity={0}
        scale={0.8}
        position={() => [0, 0]} // Center position
      />

      {/* Add a yellow Circle as the sun */}
      <Circle
        ref={sunRef}
        size={100}
        fill={'yellow'}
        position={() => [0, view.height() / 2 + 150]} // Above the beach
        opacity={0}
      />

      {/* Add a rounded yellow Rect as the beach */}
      <Rect
        ref={beachRectRef}
        width={view.width() - 20}
        height={100}
        fill={'#FFD700'} // Yellow color
        radius={20} // Rounded corners
        position={() => [0, -view.height() / 4]} // Positioned below the view
        opacity={0}
      />
    </Node>
  );

  // Animate the beach SVG's scale and opacity
  yield* all(
    beachRef().scale(1, 1.2), // Animate scale to 1 with longer duration
    beachRef().opacity(1, 1.2) // Animate opacity to 1 with longer duration
  );

  // Animate the sun's entry
  yield* sunRef().opacity(1, 1.2); // Fade the sun in
  yield* sunRef().position.y(view.height() / 2, 1.2); // Move the sun down to its position

  // Animate the beach rectangle's entry
  yield* beachRectRef().opacity(1, 1.2); // Fade in the beach
  yield* beachRectRef().position.y(view.height() / 3, 1.2); // Slide the beach up

  // Animate the fill of each Path child in the beach SVG to white
  for (const child of beachRef().children()[0].children()) {
    if (child instanceof Path) {
      yield* child.fill('white', 1.6); // Longer fill time
      yield* waitFor(0.2); // Stagger the timing for each path fill
    }
  }
});

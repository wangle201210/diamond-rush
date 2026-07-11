import java.awt.Color;
import java.awt.Font;
import java.awt.FontMetrics;
import java.awt.Graphics2D;
import java.awt.image.BufferedImage;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import javax.imageio.ImageIO;

public final class ExportFont {
    private static final int FIRST = 32;
    private static final int LAST = 126;
    private static final int COLS = 16;
    private static final int PADDING = 2;

    private ExportFont() {
    }

    public static void main(String[] args) throws Exception {
        Path output = args.length > 0 ? Path.of(args[0]) : Path.of("decoded/fonts");
        Files.createDirectories(output);
        export(output, "freej2me-small", 10, 10);
        export(output, "freej2me-medium", 12, 12);
    }

    private static void export(Path output, String name, int pointSize, int j2meHeight) throws IOException {
        Font font = new Font(Font.SANS_SERIF, Font.BOLD, pointSize);
        BufferedImage probe = new BufferedImage(1, 1, BufferedImage.TYPE_INT_ARGB);
        Graphics2D probeGraphics = probe.createGraphics();
        probeGraphics.setFont(font);
        FontMetrics metrics = probeGraphics.getFontMetrics();

        int count = LAST - FIRST + 1;
        int[] widths = new int[count];
        int cellWidth = 1;
        for (int i = 0; i < count; i++) {
            widths[i] = metrics.charWidth((char) (FIRST + i));
            cellWidth = Math.max(cellWidth, widths[i]);
        }
        int cellHeight = metrics.getHeight();
        int strideX = cellWidth + PADDING * 2;
        int strideY = cellHeight + PADDING * 2;
        int rows = (count + COLS - 1) / COLS;
        BufferedImage atlas = new BufferedImage(COLS * strideX, rows * strideY, BufferedImage.TYPE_INT_ARGB);
        Graphics2D graphics = atlas.createGraphics();
        graphics.setFont(font);
        graphics.setColor(Color.WHITE);
        for (int i = 0; i < count; i++) {
            int x = (i % COLS) * strideX + PADDING;
            int y = (i / COLS) * strideY + PADDING;
            graphics.drawString(String.valueOf((char) (FIRST + i)), x, y + metrics.getAscent() - 1);
        }
        graphics.dispose();
        probeGraphics.dispose();

        ImageIO.write(atlas, "png", output.resolve(name + ".png").toFile());
        StringBuilder json = new StringBuilder();
        json.append("{\n");
        json.append("  \"first\": ").append(FIRST).append(",\n");
        json.append("  \"last\": ").append(LAST).append(",\n");
        json.append("  \"columns\": ").append(COLS).append(",\n");
        json.append("  \"padding\": ").append(PADDING).append(",\n");
        json.append("  \"cellWidth\": ").append(cellWidth).append(",\n");
        json.append("  \"cellHeight\": ").append(cellHeight).append(",\n");
        json.append("  \"fontHeight\": ").append(j2meHeight).append(",\n");
        json.append("  \"ascent\": ").append(metrics.getAscent()).append(",\n");
        json.append("  \"widths\": [");
        for (int i = 0; i < widths.length; i++) {
            if (i > 0) {
                json.append(", ");
            }
            json.append(widths[i]);
        }
        json.append("]\n}\n");
        Files.writeString(output.resolve(name + ".json"), json.toString(), StandardCharsets.US_ASCII);
    }
}
